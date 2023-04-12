package notes

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

const (
	a4Frequency  = 440.0
	a4MidiNumber = 69
)

type Notes struct {
	config struct {
		SampleRate      *float64 `selectors:"steps"`
		SamplesPerCycle *int     `selectors:"steps"`
		DataAttributes  []string
	}
}

func (*Notes) GetName() string {
	return "notes"
}

func (*Notes) GetAllowedSelectors() []string {
	return []string{"steps", "names"}
}

func (n *Notes) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	if err := dreg.Evaluate(&n.config, pmt, slt); err != nil {
		return err
	}

	if slt.IsSelected("steps") {
		steps := make([]uint32, 0, 128)
		for note := 0; note < 128; note++ {
			freq := a4Frequency * math.Pow(2, float64(note-a4MidiNumber)/12)
			steps = append(steps, uint32((float64(*n.config.SamplesPerCycle)/(*n.config.SampleRate/freq))*(1<<16)))
		}
		r.AddData(identifier+"_steps", steps, n.config.DataAttributes, nil)
	}

	if slt.IsSelected("names") {
		prefixes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
		names := make([]string, 0, 128)
		for note := 0; note < 128; note++ {
			names = append(names, fmt.Sprintf("%s%d", prefixes[note%12], (note/12)-1))
		}
		r.AddData(identifier+"_names", names, n.config.DataAttributes, nil)
	}

	return nil
}
