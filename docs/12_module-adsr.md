# ADSR module

The ADSR module generates precomputed envelope data for Attack-Decay-Sustain-Release envelope generators. It produces envelope curve shapes, time step increments for sample-accurate envelope traversal, and human-readable description strings for UI display.

## Selectors

| Selector | Output array(s) | Type | Description |
|----------|-----------------|------|-------------|
| `curves_as3310` | `{id}_curve_as3310_attack[N]`, `{id}_curve_as3310_decay_release[N]` | 1-D | Exponential curves modeled after the CEM/AS3310 analog envelope IC |
| `curves_linear` | `{id}_curve_linear[N]` | 1-D | Linear envelope curve |
| `time_steps` | `{id}_time_steps[T]` | 1-D | Phase increment values for each time setting |
| `descriptions` | `{id}_level_descriptions[L][W]`, `{id}_time_descriptions[T][W]` | 2-D | Human-readable level and time labels |

Where `{id}` is the identifier, `N` is `adsr_samples`, `T` is `adsr_time_steps`, `L` is `adsr_level_descriptions`, and `W` is the string width.

## Parameters

| Parameter | Required by | Type | Description |
|-----------|------------|------|-------------|
| `adsr_samples` | all | `int` | Number of samples in each envelope curve |
| `adsr_sample_amplitude` | `curves_as3310`, `curves_linear` | `float64` | Peak amplitude of the envelope curve |
| `adsr_sample_scalar_type` | `curves_as3310`, `curves_linear` | `string` | C type for curve sample values (e.g., `uint8_t`) |
| `sample_rate` | `time_steps` | `float64` | Sample rate in Hz |
| `adsr_time_steps` | `time_steps`, `descriptions` | `int` | Number of discrete time settings |
| `adsr_time_steps_min_ms` | `time_steps`, `descriptions` | `int` | Minimum envelope time in milliseconds |
| `adsr_time_steps_max_ms` | `time_steps`, `descriptions` | `int` | Maximum envelope time in milliseconds |
| `adsr_time_steps_scalar_type` | `time_steps` | `string` | C type for time step values (e.g., `uint32_t`) |
| `adsr_time_steps_fractional_bit_width` | -- | `uint8` | Fractional bits for fixed-point time steps |
| `adsr_level_descriptions` | `descriptions` | `int` | Number of discrete level settings |
| `adsr_level_descriptions_string_width` | -- | `int` | Fixed string width for level labels (negative for left-aligned) |
| `adsr_time_descriptions_string_width` | -- | `int` | Fixed string width for time labels (negative for left-aligned) |
| `data_attributes` | -- | `[]string` | Optional C attributes |

## Envelope curves

### AS3310 curves

The `curves_as3310` selector generates two curves that model the behavior of the CEM/AS3310 analog envelope generator IC. The AS3310 uses an exponential charging circuit where:

- The **attack** phase charges toward a voltage higher than the peak (the IC charges toward 7.0V but clips at 5.0V). The generated attack curve replicates this by computing `1 - exp(-3t)` and scaling only the portion up to the ratio `5.0/7.0` of the asymptote. This produces the characteristically fast initial rise that gradually levels off.
- The **decay/release** phase is a standard exponential decay following `1 - exp(-3t)`, scaled to the full amplitude. The same curve is used for both decay and release since the AS3310 uses identical circuit paths for both.

Both curves are arrays of `adsr_samples` values ranging from 0 to `adsr_sample_amplitude`. The attack curve rises from 0 to the peak, and the decay/release curve is a rising ramp that firmware reverses by indexing backward.

### Linear curve

The `curves_linear` selector generates a simple linear ramp from 0 to `adsr_sample_amplitude` over `adsr_samples` values.

### Example: envelope curve usage

```c
#include "adsr-data.h"

// adsr_curve_as3310_attack is:
//   static const uint8_t adsr_curve_as3310_attack[256] = { ... };
//   #define adsr_curve_as3310_attack_len 256
//
// adsr_curve_as3310_decay_release is:
//   static const uint8_t adsr_curve_as3310_decay_release[256] = { ... };
//   #define adsr_curve_as3310_decay_release_len 256

// Read an attack envelope value.
// position: 0 to adsr_curve_as3310_attack_len - 1
uint8_t attack_value(uint8_t position) {
    return adsr_curve_as3310_attack[position];
}

// Read a decay/release envelope value (traversed in reverse).
// position: 0 to adsr_curve_as3310_decay_release_len - 1
// sustain_level: the sustain amplitude to decay toward
uint8_t decay_value(uint8_t position, uint8_t sustain_level) {
    uint8_t curve = adsr_curve_as3310_decay_release[
        adsr_curve_as3310_decay_release_len - 1 - position];
    return sustain_level + (uint16_t)curve *
        (adsr_sample_amplitude - sustain_level) / adsr_sample_amplitude;
}
```

## Time steps

The `time_steps` selector generates a 1-D array of phase increment values that control how fast the envelope traverses the curve. Each entry corresponds to a user-selectable time setting, ranging from `adsr_time_steps_min_ms` to `adsr_time_steps_max_ms`.

The time values are distributed exponentially across the range using the formula:

```
time_ms[i] = min_ms + (max_ms - min_ms) * (-1 + exp(6 * i / (steps - 1))) / (-1 + exp(6))
```

This produces a non-linear distribution where shorter times are more densely spaced (providing finer control over fast envelopes) while longer times are more spread out.

Each time step value is computed as:

```
step[i] = (adsr_samples * 1000) / (time_ms[i] * sample_rate)
```

This represents how many curve samples to advance per audio sample. When `adsr_time_steps_fractional_bit_width` is set, the value is multiplied by `2^fractional_bit_width` to produce a fixed-point number enabling sub-sample precision.

### Example: time step usage

```c
#include "adsr-data.h"

// adsr_time_steps is:
//   static const uint32_t adsr_time_steps[128] = { ... };
//   #define adsr_time_steps_len 128

// Advance the envelope phase accumulator each audio sample.
// time_setting: index into adsr_time_steps (0 = shortest, 127 = longest)
// phase_acc: fixed-point accumulator (integer part indexes into the curve)
void advance_envelope(uint8_t time_setting, uint32_t *phase_acc) {
    *phase_acc += adsr_time_steps[time_setting];
}
```

## Descriptions

The `descriptions` selector generates human-readable string arrays for displaying envelope settings on a screen or other UI.

**Level descriptions** are `adsr_level_descriptions` strings formatted as percentages from `"0.0%"` to `"100.0%"`, evenly distributed across the number of level steps.

**Time descriptions** are `adsr_time_steps` strings formatted as:
- `"Nms"` for times up to 1000 ms (e.g., `"2ms"`, `"500ms"`)
- `"N.NNs"` for times between 1000 ms and 10000 ms (e.g., `"1.50s"`, `"9.99s"`)
- `"N.Ns"` for times above 10000 ms (e.g., `"15.0s"`, `"20.0s"`)

Both arrays are emitted as 2-D `char` arrays with a fixed string width (controlled by `adsr_level_descriptions_string_width` and `adsr_time_descriptions_string_width`). A negative width produces left-aligned strings padded with spaces; a positive width produces right-aligned strings.

### Example: description usage

```c
#include "screen-data.h"

// adsr_level_descriptions is:
//   static const char adsr_level_descriptions[128][6] = { ... };
//   #define adsr_level_descriptions_rows 128
//   #define adsr_level_descriptions_cols 6
//
// adsr_time_descriptions is:
//   static const char adsr_time_descriptions[128][5] = { ... };
//   #define adsr_time_descriptions_rows 128
//   #define adsr_time_descriptions_cols 5

void display_sustain_level(uint8_t level_index) {
    // level_index ranges from 0 to adsr_level_descriptions_rows - 1
    draw_text(adsr_level_descriptions[level_index]);  // e.g., "50.4%"
}

void display_attack_time(uint8_t time_index) {
    // time_index ranges from 0 to adsr_time_descriptions_rows - 1
    draw_text(adsr_time_descriptions[time_index]);  // e.g., "2ms" or "1.50s"
}
```

## Example configuration

```yaml
global_parameters:
  sample_rate: 48000
  adsr_samples: 0x0100
  adsr_sample_amplitude: 0xff
  adsr_sample_scalar_type: uint8_t
  adsr_time_steps: 0x80
  adsr_time_steps_min_ms: 2
  adsr_time_steps_max_ms: 20000
  adsr_time_steps_scalar_type: uint32_t
  adsr_time_steps_fractional_bit_width: 16
  adsr_level_descriptions: 0x80
  adsr_level_descriptions_string_width: -6
  adsr_time_descriptions_string_width: -5

output:
  firmware/adsr-data.h:
    includes:
      stdint.h: true
    macros:
      adsr_sample_amplitude:
        value: 0xff
        type: uint8_t
        hex: true
    modules:
      adsr:
        name: adsr
        selectors:
          - curves_as3310
          - curves_linear
          - time_steps

  firmware/screen-data.h:
    modules:
      adsr:
        name: adsr
        selectors:
          - descriptions
```
