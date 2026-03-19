# Filters module

The filters module generates precomputed filter coefficient tables for digital audio filters. It currently supports one-pole (first-order IIR) low-pass and high-pass filters, along with human-readable frequency description strings for UI display.

## Selectors

| Selector | Output array | Type | Description |
|----------|-------------|------|-------------|
| `lowpass_onepole` | `{id}_lowpass_onepole_coefficients[N]` | 1-D (struct) | Low-pass one-pole filter coefficients |
| `highpass_onepole` | `{id}_highpass_onepole_coefficients[N]` | 1-D (struct) | High-pass one-pole filter coefficients |
| `descriptions` | `{id}_frequency_descriptions[N][W]` | 2-D (char) | Human-readable frequency labels |

Where `{id}` is the identifier, `N` is `filters_frequencies`, and `W` is the string width.

## Parameters

| Parameter | Required by | Type | Description |
|-----------|------------|------|-------------|
| `sample_rate` | all | `float64` | Sample rate in Hz |
| `filters_frequencies` | all | `int` | Number of discrete cutoff frequency settings |
| `filters_frequency_min` | all | `float64` | Minimum cutoff frequency in Hz |
| `filters_frequency_max` | all | `float64` | Maximum cutoff frequency in Hz |
| `filters_coefficients_onepole_scalar_type` | `lowpass_onepole`, `highpass_onepole` | `string` | C type for coefficient values (e.g., `int8_t`) |
| `filters_coefficients_onepole_fractional_bit_width` | -- | `uint8` | Fractional bits for fixed-point coefficients |
| `filters_frequency_descriptions_string_width` | -- | `int` | Fixed string width for frequency labels (negative for left-aligned) |
| `data_attributes` | -- | `[]string` | Optional C attributes |

## Frequency distribution

The cutoff frequencies are distributed exponentially between `filters_frequency_min` and `filters_frequency_max` using the formula:

```
freq[i] = freq_min + (freq_max - freq_min) * (-1 + exp(3 * i / (N - 1))) / (-1 + exp(3))
```

This produces a perceptually useful distribution where lower frequencies are more densely spaced (important for musical filter sweeps) while higher frequencies are more spread out. The exponential base of 3 gives approximately a 20:1 density ratio between the low end and high end of the range.

## One-pole filter coefficients

The one-pole filters implement first-order IIR filters using the bilinear transform. Each entry in the coefficient array is a struct with three fields: `a1`, `b0`, and `b1`.

For each cutoff frequency `fc`, the normalized frequency is `fn = fc / sample_rate`, and the coefficient `a1` is computed as:

```
a1 = (1 - tan(pi * fn)) / (1 + tan(pi * fn))
```

### Low-pass coefficients

For the low-pass filter:

```
b0 = (1 - a1) / 2
b1 = (1 - a1) / 2
```

The difference equation is: `y[n] = b0 * x[n] + b1 * x[n-1] + a1 * y[n-1]`

### High-pass coefficients

For the high-pass filter:

```
b0 = (1 + a1) / 2
b1 = -(1 + a1) / 2
```

The difference equation is the same: `y[n] = b0 * x[n] + b1 * x[n-1] + a1 * y[n-1]`

### Fixed-point scaling

When `filters_coefficients_onepole_fractional_bit_width` is set, all three coefficients (`a1`, `b0`, `b1`) are multiplied by `2^fractional_bit_width` before conversion to the target integer type. Firmware must right-shift the intermediate products by the same number of bits after multiplication.

### Generated struct format

The coefficient arrays are emitted as arrays of anonymous C structs:

```c
static const struct {
    int8_t a1;
    int8_t b0;
    int8_t b1;
} filter_lowpass_onepole_coefficients[128] = {
    {a1_0, b0_0, b1_0},
    {a1_1, b0_1, b1_1},
    ...
};
#define filter_lowpass_onepole_coefficients_len 128
```

### Example: filter usage

```c
#include "filter-data.h"

// filter_lowpass_onepole_coefficients is:
//   static const struct { int8_t a1; int8_t b0; int8_t b1; }
//       filter_lowpass_onepole_coefficients[128] = { ... };
//   #define filter_lowpass_onepole_coefficients_len 128

// Apply one-pole low-pass filter to a sample.
// cutoff_index: index into the coefficient table (0 to 127)
// x_n: current input sample
// x_prev: previous input sample
// y_prev: previous output sample
// Returns the filtered output sample.
// Assumes 7-bit fractional coefficients (fractional_bit_width = 7).
int16_t lowpass_filter(uint8_t cutoff_index, int16_t x_n,
                       int16_t x_prev, int16_t y_prev) {
    int8_t a1 = filter_lowpass_onepole_coefficients[cutoff_index].a1;
    int8_t b0 = filter_lowpass_onepole_coefficients[cutoff_index].b0;
    int8_t b1 = filter_lowpass_onepole_coefficients[cutoff_index].b1;

    int16_t result = ((int32_t)b0 * x_n +
                      (int32_t)b1 * x_prev +
                      (int32_t)a1 * y_prev) >> 7;
    return result;
}
```

## Descriptions

The `descriptions` selector generates frequency labels for each cutoff frequency setting:

- Frequencies at or below 1000 Hz are formatted as `"NHz"` (e.g., `"20Hz"`, `"500Hz"`)
- Frequencies above 1000 Hz are formatted as `"N.NNkHz"` (e.g., `"1.50kHz"`, `"20.00kHz"`)

The strings are emitted as a 2-D `char` array with a fixed width controlled by `filters_frequency_descriptions_string_width`. A negative width produces left-aligned strings.

### Example: description usage

```c
#include "screen-data.h"

// filter_frequency_descriptions is:
//   static const char filter_frequency_descriptions[128][8] = { ... };
//   #define filter_frequency_descriptions_rows 128
//   #define filter_frequency_descriptions_cols 8

void display_filter_frequency(uint8_t freq_index) {
    draw_text(filter_frequency_descriptions[freq_index]);  // e.g., "1.50kHz"
}
```

## Example configuration

```yaml
global_parameters:
  sample_rate: 48000
  filters_frequencies: 0x80
  filters_frequency_min: 20
  filters_frequency_max: 20000
  filters_frequency_descriptions_string_width: -8
  filters_coefficients_onepole_scalar_type: int8_t
  filters_coefficients_onepole_fractional_bit_width: 7

output:
  firmware/filter-data.h:
    includes:
      stdint.h: true
    modules:
      filter:
        name: filters
        selectors:
          - lowpass_onepole
          - highpass_onepole

  firmware/screen-data.h:
    modules:
      filter:
        name: filters
        selectors:
          - descriptions
```
