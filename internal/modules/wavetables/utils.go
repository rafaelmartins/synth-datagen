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
	return math.Sqrt(noteFrequency(min(octave*12, 127), a4Freq) * noteFrequency(min((octave+1)*12-1, 127), a4Freq))
}
