# DSP modules

synth-datagen includes four DSP modules that generate precomputed data arrays for synthesizer firmware. Each module is invoked from the YAML configuration file by name, with selectors controlling which specific data arrays to produce. The generated C headers declare all data as `static const` arrays with accompanying `#define` macros for dimensions.

## Module overview

| Module | Name in config | Purpose |
|--------|---------------|---------|
| [Wavetables](11_module-wavetables.md) | `wavetables` | Waveform lookup tables (sine, square, triangle, sawtooth) and band-limited variants |
| [ADSR](12_module-adsr.md) | `adsr` | Envelope curve shapes, time step increments, and human-readable descriptions |
| [Filters](13_module-filters.md) | `filters` | One-pole low-pass and high-pass filter coefficients, and human-readable descriptions |
| [Notes](14_module-notes.md) | `notes` | MIDI note phase steps, note names, and octave numbers |

## Common patterns

### Selectors

Each module defines a set of allowed selectors. The configuration file lists which selectors to activate for a given output, controlling which data arrays appear in the generated header. For example, the wavetables module supports `sine`, `square`, `triangle`, `sawtooth`, `blsquare`, `bltriangle`, and `blsawtooth`. Selecting `sine` and `blsquare` produces only those two data arrays.

### Identifier prefixing

The `identifier` key in the configuration determines the C identifier prefix for all data produced by that module invocation. For instance, if the identifier is `oscillator` and the selector is `sine`, the generated array is named `oscillator_sine`.

### Dimension macros

The code generator automatically emits dimension macros alongside each array:

| Array dimensions | Generated macros |
|-----------------|-----------------|
| 1-D array (`type name[N]`) | `#define name_len N` |
| 2-D array (`type name[R][C]`) | `#define name_rows R` and `#define name_cols C` |
| 3+ dimensions (`type name[A][B][C]`) | `#define name_len_0 A`, `#define name_len_1 B`, `#define name_len_2 C` |

### Data attributes

All modules support a `data_attributes` parameter (passed via `parameters` in the config). When set, the attribute strings are inserted between the variable name and the initializer in the generated C declaration. This is commonly used for AVR `PROGMEM`:

```c
static const int16_t oscillator_sine[512] PROGMEM = { ... };
```

### Scalar types and fixed-point arithmetic

Several modules support configurable scalar types and fractional bit widths via `*_scalar_type` and `*_fractional_bit_width` parameters. The tool supports the full range of C scalar types:

- **Integer types** (`uint8_t`, `int16_t`, `uint32_t`, etc.) -- for fixed-point arithmetic on systems without an FPU. When a fractional bit width is specified, values are multiplied by `2^fractional_bit_width` before conversion to the integer type, producing fixed-point values that firmware can process with integer arithmetic and bit-shift operations.
- **Single-precision floating-point** (`float`) -- for platforms equipped with a single-precision FPU. Setting the scalar type to `float` produces `float` arrays directly usable without any fixed-point conversion in firmware.
- **Double-precision floating-point** (`double`) -- for platforms equipped with a double-precision FPU. Setting the scalar type to `double` produces `double` arrays with full double-precision values.

When using `float` or `double`, the generated data contains natural floating-point values that can be used directly in arithmetic without any scaling or bit-shifting.

> [!WARNING]
> The `*_fractional_bit_width` parameters remain functional even when using `float` or `double` scalar types. When set, the computed values are multiplied by `2^fractional_bit_width` before being stored as floats, producing scaled values rather than natural ones. This can be intentional when simulating fixed-point quantization behavior on a floating-point platform, but users targeting full floating-point precision should omit the fractional bit width parameters entirely to avoid unexpected scaling.

## Configuration structure

Each module invocation in the YAML configuration follows this pattern:

```yaml
output:
  firmware/my-data.h:
    modules:
      my_identifier:          # C identifier prefix
        name: wavetables      # module name
        selectors:            # which data arrays to generate
          - sine
          - blsquare
        parameters:           # optional per-invocation overrides
          data_attributes:
            - PROGMEM
```

### Parameter resolution

Each module declares an internal configuration struct whose fields are populated automatically by the data registry. The resolution process for each field works as follows:

1. Convert the Go field name to snake\_case (e.g., `SampleAmplitude` → `sample_amplitude`)
2. Search the per-invocation `parameters` map for `{module}_{field}` (e.g., `wavetables_sample_amplitude`), then `{field}` (e.g., `sample_amplitude`)
3. If not found locally, repeat the same two-key search in `global_parameters`
4. If still not found: pointer and slice fields remain nil (optional); non-pointer, non-slice fields cause an error (required)

Some struct fields carry a `selectors` tag (e.g., `SampleRate` with tag `selectors:"blsquare,bltriangle,blsawtooth"`). These fields are only required when at least one of the listed selectors is active. If none of the listed selectors are active and the field is a pointer, it is left nil without error.

This design allows a single `global_parameters` block to supply defaults to all module invocations, while individual invocations can override specific values through their `parameters` map. See [Configuration](20_configuration.md) for the full resolution rules and examples.
