package wavetables

import (
	"math"
)

const (
	a4Frequency  = 440.0
	a4MidiNumber = 69
)

func noteFrequency(note int, a4Freq float64) float64 {
	return a4Freq * math.Pow(2, float64(note-a4MidiNumber)/12)
}

func wavetableFrequency(octave int, a4Freq float64) float64 {
	first := octave * 12
	if first > 127 {
		first = 127
	}
	last := (octave+1)*12 - 1
	if last > 127 {
		last = 127
	}

	return math.Sqrt(noteFrequency(first, a4Freq) * noteFrequency(last, a4Freq))
}
