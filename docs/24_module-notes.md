# Notes module

The notes module generates MIDI note-related lookup tables: phase step values for oscillator frequency control, note name strings, and octave number mappings. All tables cover the full 128-note MIDI range (notes 0 through 127).

## Selectors

| Selector | Output array | Type | Description |
|----------|-------------|------|-------------|
| `phase_steps` | `{id}_phase_steps[128]` | 1-D | Phase increment per audio sample for each MIDI note |
| `names` | `{id}_names[128]` | 1-D (string) | Note name strings (e.g., `"C4"`, `"A#5"`) |
| `octaves` | `{id}_octaves[128]` | 1-D (`uint8_t`) | Octave number for each MIDI note (0 through 10) |

Where `{id}` is the identifier from the configuration.

## Parameters

| Parameter | Required by | Type | Description |
|-----------|------------|------|-------------|
| `a4_frequency` | -- | `float64` | Reference frequency for A4 (defaults to 440.0 Hz) |
| `sample_rate` | `phase_steps` | `float64` | Sample rate in Hz |
| `samples_per_cycle` | `phase_steps` | `int` | Number of samples in one waveform cycle |
| `notes_phase_steps_scalar_type` | `phase_steps` | `string` | C type for phase step values (e.g., `uint32_t`) |
| `notes_phase_steps_fractional_bit_width` | -- | `uint8` | Fractional bits for fixed-point phase steps |
| `data_attributes` | -- | `[]string` | Optional C attributes |

## Phase steps

The `phase_steps` selector generates a 128-element array where each entry is the phase increment to add to a phase accumulator for each audio sample to produce the corresponding MIDI note's frequency.

For each MIDI note number `n`, the frequency is computed using equal temperament:

```
freq[n] = a4_frequency * 2^((n - 69) / 12)
```

Where 69 is the MIDI note number for A4. The phase step is then:

```
step[n] = samples_per_cycle * freq[n] / sample_rate
```

This value represents how many wavetable samples to advance per audio sample. For a 512-sample wavetable at 48000 Hz sample rate, A4 (440 Hz) produces a phase step of approximately `512 * 440 / 48000 ≈ 4.693`.

When `notes_phase_steps_fractional_bit_width` is set, the values are multiplied by `2^fractional_bit_width` to produce fixed-point integers. With a 16-bit fractional width, the A4 phase step becomes approximately `4.693 * 65536 ≈ 307,626`.

### Example: phase step usage

```c
#include "oscillator-data.h"

// notes_phase_steps is:
//   static const uint32_t notes_phase_steps[128] = { ... };
//   #define notes_phase_steps_len 128
//
// oscillator_sine is:
//   static const int16_t oscillator_sine[512] PROGMEM = { ... };
//   #define oscillator_sine_len 512

// Phase accumulator for the oscillator.
// High bits index into the wavetable, low bits are the fractional part.
static uint32_t phase_accumulator = 0;

// Advance the oscillator and return the current sample.
// midi_note: MIDI note number (0 to 127)
// Assumes 16-bit fractional phase steps.
int16_t oscillator_tick(uint8_t midi_note) {
    phase_accumulator += notes_phase_steps[midi_note];
    uint16_t index = (phase_accumulator >> 16) % oscillator_sine_len;
    return oscillator_sine[index];
}
```

## Note names

The `names` selector generates a 128-element array of note name strings following the standard naming convention: `C`, `C#`, `D`, `D#`, `E`, `F`, `F#`, `G`, `G#`, `A`, `A#`, `B`, combined with an octave number starting at -1.

| MIDI note | Name |
|-----------|------|
| 0 | `C-1` |
| 12 | `C0` |
| 60 | `C4` |
| 69 | `A4` |
| 127 | `G9` |

The names are emitted as a `char*` string array:

```c
static const char* notes_names[128] = {
    "C-1", "C#-1", "D-1", ..., "G9",
};
#define notes_names_len 128
```

## Octave numbers

The `octaves` selector generates a 128-element `uint8_t` array where each entry is the octave number for the corresponding MIDI note, computed as `note / 12`. This produces values from 0 (MIDI notes 0--11) to 10 (MIDI note 120--127).

> [!NOTE]
> The octave numbering here starts at 0 (for MIDI notes 0--11), not at -1 as in the note names. This is because the octave values map directly to band-limited wavetable rows, where octave 0 corresponds to the first row.

```c
static const uint8_t notes_octaves[128] = {
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,  // MIDI 0-11
    1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,  // MIDI 12-23
    ...
};
#define notes_octaves_len 128
```

### Example: octave lookup for band-limited wavetables

```c
#include "oscillator-data.h"

// notes_octaves is:
//   static const uint8_t notes_octaves[128] = { ... };
//   #define notes_octaves_len 128
//
// oscillator_blsquare is:
//   static const int16_t oscillator_blsquare[10][512] PROGMEM = { ... };
//   #define oscillator_blsquare_rows 10
//   #define oscillator_blsquare_cols 512

int16_t read_blsquare(uint8_t midi_note, uint16_t phase) {
    uint8_t octave = notes_octaves[midi_note];
    if (octave >= oscillator_blsquare_rows)
        octave = oscillator_blsquare_rows - 1;
    uint16_t index = (phase >> 7) % oscillator_blsquare_cols;
    return oscillator_blsquare[octave][index];
}
```

## Example configuration

```yaml
global_parameters:
  sample_rate: 48000
  samples_per_cycle: 0x0200
  notes_phase_steps_scalar_type: uint32_t
  notes_phase_steps_fractional_bit_width: 16

output:
  firmware/oscillator-data.h:
    includes:
      stdint.h: true
    modules:
      notes:
        name: notes
        selectors:
          - phase_steps
          - octaves
```
