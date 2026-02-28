# Configuration

synth-datagen reads a YAML configuration file (default: `synth-datagen.yml`) that defines global parameters and one or more output header files. This page documents the configuration file format, command-line options, and all supported configuration keys.

## Command-line usage

```bash
synth-datagen [-f config.yml] [-o output_dir] [-c] [-v]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | `synth-datagen.yml` | Path to the configuration file |
| `-o` | `.` | Output directory for generated files |
| `-c` | disabled | Generate HTML chart files instead of C headers |
| `-v` | -- | Print version and exit |

When `-c` is specified, only outputs with a `charts_output` field are processed, and the tool generates HTML visualization files (using go-echarts) instead of C headers.

## Configuration file structure

The top-level YAML structure has two keys:

```yaml
global_parameters:
  # key-value pairs accessible to all modules

output:
  # mapping of output file paths to their content definitions
```

### Global parameters

The `global_parameters` section is a flat key-value map that provides default values to all module invocations. Each DSP module internally defines a configuration struct with typed fields (e.g., the ADSR module has fields `Samples`, `SampleAmplitude`, `SampleScalarType`). The parameter resolver (data registry) populates these fields by converting each field name to snake\_case and searching two maps in order:

1. The per-invocation `parameters` map (local override)
2. The `global_parameters` map (global default)

Within each map, two key forms are checked, in order:

1. **Module-prefixed key**: `{module}_{field}` (e.g., `adsr_sample_amplitude`)
2. **Unprefixed key**: `{field}` (e.g., `sample_amplitude`)

For example, the ADSR module's `SampleAmplitude` field (snake\_case: `sample_amplitude`) is resolved for a module named `adsr` as follows:

1. Check local `parameters` for `adsr_sample_amplitude` → if found, use it
2. Check local `parameters` for `sample_amplitude` → if found, use it
3. Check `global_parameters` for `adsr_sample_amplitude` → if found, use it
4. Check `global_parameters` for `sample_amplitude` → if found, use it
5. If still not found: pointer and slice fields remain nil (optional); non-pointer, non-slice fields cause an error (required)

This means `adsr_samples: 0x0100` in `global_parameters` matches the ADSR module's `Samples` field because the module-prefixed lookup `adsr_` + `samples` = `adsr_samples` succeeds. A shared parameter like `sample_rate: 48000` (unprefixed) is found as a fallback for any module that needs a `SampleRate` field.

To override a global parameter for a specific module invocation, add it to the `parameters` map of that module:

```yaml
global_parameters:
  sample_rate: &sample_rate 48000
  wavetables_sample_amplitude: 0x01ff

output:
  firmware/oscillator-data.h:
    modules:
      oscillator:
        name: wavetables
        selectors:
          - sine
        parameters:
          # overrides global wavetables_sample_amplitude for this invocation only
          wavetables_sample_amplitude: 0x00ff
```

YAML anchors and aliases can be used to reference values across the configuration:

```yaml
global_parameters:
  sample_rate: &sample_rate 48000
  wavetables_sample_amplitude: &wavetables_sample_amplitude 0x01ff

output:
  firmware/main-data.h:
    macros:
      sample_rate: *sample_rate
```

### Output definitions

Each key under `output` is the relative path for the generated header file. Each output contains:

| Field | Type | Description |
|-------|------|-------------|
| `charts_output` | `string` | Optional path for HTML chart output (used with `-c` flag) |
| `includes` | mapping | C `#include` directives |
| `macros` | mapping | C `#define` preprocessor macros |
| `variables` | mapping | C `static const` variable declarations |
| `modules` | mapping | DSP module invocations |

## Includes

The `includes` section maps header file paths to a boolean indicating whether they are system includes:

```yaml
includes:
  stdint.h: true       # generates: #include <stdint.h>
  avr/pgmspace.h: true # generates: #include <avr/pgmspace.h>
  config.h: false      # generates: #include "config.h"
```

Duplicate includes are deduplicated. If the same path appears with both `true` (system) and `false` (local), the local form takes precedence.

## Macros

Each key under `macros` becomes the `#define` identifier. Values can be specified as plain scalars or as objects with additional options:

### Simple scalar macro

```yaml
macros:
  sample_rate: 48000
```

Generates: `#define sample_rate 48000`

### Typed macro

```yaml
macros:
  adsr_sample_amplitude:
    value: 0xff
    type: uint8_t
    hex: true
```

Generates: `#define adsr_sample_amplitude 0xff`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | any scalar | -- | The macro value |
| `type` | `string` | -- | C type for value conversion (affects hex formatting width) |
| `hex` | `bool` | `false` | Format numeric values in hexadecimal |
| `raw` | `bool` | `false` | Emit the value as-is without type formatting |

### Raw macro

When `raw` is `true`, the value is emitted verbatim without any type conversion or formatting. This is useful for referencing platform-specific constants:

```yaml
macros:
  cpu_frqsel:
    value: CLKCTRL_FRQSEL_24M_gc
    raw: true
```

Generates: `#define cpu_frqsel CLKCTRL_FRQSEL_24M_gc`

### Expression evaluation

When `eval_env` is provided (or `eval` is `true`), string values are evaluated as expressions using the [expr](https://github.com/expr-lang/expr) library. The `eval_env` map provides variables available in the expression:

```yaml
macros:
  timer_tcb_ccmp:
    value: clock_frequency / sample_rate
    type: uint16_t
    eval_env:
      clock_frequency: 24000000
      sample_rate: 48000
```

Generates: `#define timer_tcb_ccmp 500`

This is useful for computing derived constants from base parameters, such as baud rate register values or timer periods.

## Variables

The `variables` section works similarly to macros but generates `static const` variable declarations instead of `#define` macros. Each key becomes the variable identifier:

```yaml
variables:
  my_data:
    value: [1, 2, 3, 4]
    type: uint8_t
```

Generates:

```c
static const uint8_t my_data[4] = {
    0x01, 0x02, 0x03, 0x04,
};
#define my_data_len 4
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `value` | scalar or array | -- | The variable's value |
| `type` | `string` | -- | C type for the value |
| `string_width` | `int` | -- | Fixed string width (negative for left-aligned) |
| `attributes` | `[]string` | -- | C attributes inserted before the initializer |
| `eval` | `bool` | `false` | Evaluate string values as expressions |
| `eval_env` | mapping | -- | Variables for expression evaluation |

## Modules

The `modules` section invokes DSP modules to generate data arrays. Each key becomes the C identifier prefix:

```yaml
modules:
  oscillator:
    name: wavetables
    selectors:
      - sine
      - blsquare
    parameters:
      data_attributes:
        - PROGMEM
```

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Module name (`wavetables`, `adsr`, `filters`, or `notes`) |
| `selectors` | `[]string` | Which data arrays to generate |
| `parameters` | mapping | Per-invocation parameter overrides |

The `parameters` map is checked before `global_parameters` during parameter resolution. See [DSP modules](10_modules.md) for detailed documentation of each module's selectors and parameters.

## Supported C types

The following C scalar types are supported for `type` fields and module `*_scalar_type` parameters:

| C type | Go representation | Notes |
|--------|------------------|-------|
| `bool` | `bool` | |
| `int8_t` | `int8` | |
| `int16_t` | `int16` | |
| `int32_t` | `int32` | |
| `int64_t` | `int64` | |
| `uint8_t` | `uint8` | |
| `uint16_t` | `uint16` | |
| `uint32_t` | `uint32` | |
| `uint64_t` | `uint64` | |
| `float` | `float32` | Single-precision FPU targets |
| `double` | `float64` | Double-precision FPU targets |
| `char*` | `string` | |

Setting a module's `*_scalar_type` parameter to `float` or `double` produces floating-point arrays that can be used directly on platforms with an FPU, without any fixed-point scaling in firmware. See [DSP modules -- Scalar types and fixed-point arithmetic](10_modules.md) for details on how fractional bit width parameters interact with floating-point types.

## Complete example

The following is a minimal but complete configuration generating wavetable and note data:

```yaml
global_parameters:
  sample_rate: &sample_rate 48000
  samples_per_cycle: 0x0200
  wavetables_sample_amplitude: 0x01ff
  wavetables_sample_scalar_type: int16_t
  notes_phase_steps_scalar_type: uint32_t
  notes_phase_steps_fractional_bit_width: 16

output:
  firmware/oscillator-data.h:
    includes:
      stdint.h: true
    modules:
      oscillator:
        name: wavetables
        selectors:
          - sine
      notes:
        name: notes
        selectors:
          - phase_steps
```

Running `synth-datagen -f synth-datagen.yml -o .` generates `firmware/oscillator-data.h` containing a 512-element `int16_t` sine wavetable and a 128-element `uint32_t` phase step table.
