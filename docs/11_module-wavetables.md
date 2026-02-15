# Wavetables module

The wavetables module generates waveform lookup tables for oscillator implementations. It produces both naive (mathematically ideal) waveforms and band-limited variants suitable for alias-free synthesis.

## Selectors

| Selector | Output array | Type | Description |
|----------|-------------|------|-------------|
| `sine` | `{id}_sine[N]` | 1-D | One cycle of a sine wave |
| `square` | `{id}_square[N]` | 1-D | One cycle of a naive square wave |
| `triangle` | `{id}_triangle[N]` | 1-D | One cycle of a naive triangle wave |
| `sawtooth` | `{id}_sawtooth[N]` | 1-D | One cycle of a naive sawtooth wave |
| `blsquare` | `{id}_blsquare[R][C]` | 2-D | Band-limited square wave, one row per octave |
| `bltriangle` | `{id}_bltriangle[R][C]` | 2-D | Band-limited triangle wave, one row per octave |
| `blsawtooth` | `{id}_blsawtooth[R][C]` | 2-D | Band-limited sawtooth wave, one row per octave |

Where `{id}` is the identifier from the configuration, `N` is `samples_per_cycle`, `R` is the number of octaves, and `C` is `samples_per_cycle`.

## Parameters

| Parameter | Required by | Type | Description |
|-----------|------------|------|-------------|
| `samples_per_cycle` | all | `int` | Number of samples in one waveform cycle |
| `wavetables_sample_amplitude` | all | `float64` | Peak amplitude of the waveform |
| `wavetables_sample_scalar_type` | all | `string` | C type for sample values (e.g., `int16_t`) |
| `data_attributes` | -- | `[]string` | Optional C attributes (e.g., `PROGMEM`) |
| `sample_rate` | `blsquare`, `bltriangle`, `blsawtooth` | `float64` | Sample rate in Hz, used to compute harmonics |
| `a4_frequency` | -- | `float64` | Reference frequency for A4 (defaults to 440.0 Hz) |
| `wavetables_bandlimited_omit_high_octaves` | -- | `int` | Number of highest octaves to exclude from band-limited tables |

Parameters are resolved from `global_parameters` or per-module `parameters` overrides. The parameter resolver checks for `wavetables_`-prefixed keys first, then falls back to unprefixed keys.

## Naive waveforms

The naive waveforms compute one cycle of `samples_per_cycle` samples. Each sample is computed mathematically and scaled to the range `[-sample_amplitude, +sample_amplitude]`.

**Sine** computes `amplitude * sin(2 * pi * i / samples_per_cycle)` for each sample index `i`.

**Square** computes `amplitude * (1 - 2 * floor(2 * i / samples_per_cycle))`, producing a waveform that transitions between `+amplitude` and `-amplitude` at the midpoint.

**Triangle** computes `amplitude * (2/pi) * asin(sin(2 * pi * i / samples_per_cycle))`, producing a linear ramp between peaks.

**Sawtooth** computes `amplitude * (1 - 2 * i / samples_per_cycle)`, producing a linear ramp from `+amplitude` to `-amplitude`.

### Example: naive waveform usage

Given a configuration producing `oscillator_sine` with `samples_per_cycle = 0x0200` (512) and `int16_t` scalar type:

```c
#include "oscillator-data.h"

// oscillator_sine is:
//   static const int16_t oscillator_sine[512] PROGMEM = { ... };
//   #define oscillator_sine_len 512

int16_t read_sine_sample(uint16_t phase) {
    uint16_t index = phase >> 7;  // map 16-bit phase to 9-bit index
    return oscillator_sine[index % oscillator_sine_len];
}
```

## Band-limited waveforms

The band-limited waveforms use BLIT (Band-Limited Impulse Train) synthesis to produce alias-free waveform tables. A separate table is generated for each MIDI octave (128 MIDI notes / 12 notes per octave = 11 octaves, ceiling), minus any octaves excluded by `bandlimited_omit_high_octaves`.

For each octave, the algorithm:

1. Computes a representative frequency for the octave as the geometric mean of the lowest and highest note frequencies in that octave
2. Determines the number of harmonics that fit below the Nyquist frequency (`sample_rate / 2`)
3. Generates a BLIT kernel: `sin(pi * x * harmonics) / (harmonics * sin(pi * x))` over one cycle
4. Integrates the BLIT to produce a square wave (for `blsquare` and as an intermediate for `bltriangle`)
5. Integrates the square wave to produce a triangle wave (for `bltriangle`), with a quarter-cycle phase shift
6. Integrates an offset BLIT to produce a sawtooth wave (for `blsawtooth`)
7. Normalizes each resulting waveform to the full `[-sample_amplitude, +sample_amplitude]` range

The result is a 2-D array where each row contains the waveform for one octave. At higher octaves, fewer harmonics are included, progressively smoothing the waveform to avoid aliasing.

### Example: band-limited waveform usage

Given a configuration producing `oscillator_blsquare` with 10 octaves and 512 samples per cycle:

```c
#include "oscillator-data.h"

// oscillator_blsquare is:
//   static const int16_t oscillator_blsquare[10][512] PROGMEM = { ... };
//   #define oscillator_blsquare_rows 10
//   #define oscillator_blsquare_cols 512

int16_t read_blsquare_sample(uint8_t midi_note, uint16_t phase) {
    uint8_t octave = midi_note / 12;
    if (octave >= oscillator_blsquare_rows)
        octave = oscillator_blsquare_rows - 1;
    uint16_t index = phase >> 7;  // map phase to table index
    return oscillator_blsquare[octave][index % oscillator_blsquare_cols];
}
```

## Example configuration

```yaml
global_parameters:
  sample_rate: 48000
  samples_per_cycle: 0x0200
  wavetables_sample_amplitude: 0x01ff
  wavetables_sample_scalar_type: int16_t
  wavetables_bandlimited_omit_high_octaves: 1

output:
  firmware/oscillator-data.h:
    includes:
      avr/pgmspace.h: true
      stdint.h: true
    modules:
      oscillator:
        name: wavetables
        selectors:
          - sine
          - blsquare
          - bltriangle
          - blsawtooth
        parameters:
          data_attributes:
            - PROGMEM
```

This configuration generates `oscillator_sine` as a 1-D array of 512 `int16_t` values with amplitude range `[-511, +511]`, and `oscillator_blsquare`, `oscillator_bltriangle`, and `oscillator_blsawtooth` as 2-D arrays with 10 octave rows (11 minus 1 omitted) of 512 samples each. All arrays carry the `PROGMEM` attribute for AVR flash storage.
