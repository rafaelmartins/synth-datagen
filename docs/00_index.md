---
menu: Main
---
**A tool that generates C data headers for synthesizer waveforms and DSP algorithms.**

## Overview

synth-datagen is a command-line tool written in Go that generates precomputed C header files for use in synthesizer firmware. It reads a YAML configuration file describing the desired output and produces `static const` data arrays, preprocessor macros, and struct definitions ready to be `#include`d by C code.

The tool supports a range of target platforms -- from resource-constrained microcontrollers (such as AVR) that require fixed-point integer arithmetic and flash storage attributes like `PROGMEM`, to systems equipped with single-precision or double-precision FPUs where native `float` or `double` arrays can be used directly. The scalar type for each data set is configurable, so the same configuration structure can target different hardware by changing a type parameter.

## Key highlights

- **Precomputed DSP data** -- generates wavetables, ADSR envelope curves, filter coefficients, and MIDI note phase steps as C arrays
- **Band-limited waveforms** -- produces per-octave band-limited square, triangle, and sawtooth tables using BLIT synthesis to avoid aliasing
- **Fixed-point and floating-point support** -- configurable scalar types from `uint8_t` to `double`, with optional fractional bit widths for integer-based fixed-point arithmetic or native `float`/`double` output for FPU-equipped platforms
- **Flexible output** -- each output header file independently selects which modules and selectors to include, with per-output macros and includes
- **Expression evaluation** -- macro and variable values can be computed from expressions with custom environments, enabling derived constants like baud rate registers
- **Chart generation** -- optional HTML chart output for visual inspection of generated waveforms and curves

## Explore further

- [Getting started](10_getting-started.md) -- installation and overview of the tool's workflow
- [DSP modules](20_modules.md) -- detailed documentation for each DSP module and the C arrays it produces
- [Configuration](30_configuration.md) -- YAML configuration file format, global parameters, macros, and output structure
- [Source code](https://github.com/rafaelmartins/synth-datagen) -- GitHub repository
