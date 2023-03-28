package blwavetables

import (
	"math"
)

const (
	a4Frequency  = 440.0
	a4MidiNumber = 69
)

func noteFrequency(note int) float64 {
	return a4Frequency * math.Pow(2, float64(note-a4MidiNumber)/12)
}

func wavetableFrequency(octave int) float64 {
	first := octave * 12
	if first > 127 {
		first = 127
	}
	last := (octave+1)*12 - 1
	if last > 127 {
		last = 127
	}

	return math.Sqrt(noteFrequency(first) * noteFrequency(last))
}
