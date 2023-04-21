package notes

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

const (
	a4Frequency  = 440.0
	a4MidiNumber = 69
)

type Notes struct{}

type notesConfig struct {
	DataAttributes               []string
	A4Frequency                  *float64
	SampleRate                   *float64 `selectors:"phase_steps"`
	SamplesPerCycle              *int     `selectors:"phase_steps"`
	PhaseStepsScalarType         *string  `selectors:"phase_steps"`
	PhaseStepsFractionalBitWidth *uint8   `selectors:"phase_steps"`
}

func (*Notes) GetName() string {
	return "notes"
}

func (*Notes) GetAllowedSelectors() []string {
	return []string{"phase_steps", "names", "octaves"}
}

func (n *Notes) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	config := notesConfig{}
	if err := dreg.Evaluate(n.GetName(), &config, pmt, slt); err != nil {
		return err
	}

	if config.A4Frequency == nil {
		config.A4Frequency = new(float64)
		*config.A4Frequency = a4Frequency
	}

	if slt.IsSelected("phase_steps") {
		steps := make([]uint64, 0, 128)
		for note := 0; note < 128; note++ {
			freq := *config.A4Frequency * math.Pow(2, float64(note-a4MidiNumber)/12)
			steps = append(steps, uint64((float64(*config.SamplesPerCycle)/(*config.SampleRate/freq))*float64(int(1)<<*config.PhaseStepsFractionalBitWidth)))
		}

		s, err := convert.Slice(steps, *config.PhaseStepsScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_phase_steps", s, config.DataAttributes, nil)
	}

	if slt.IsSelected("names") {
		prefixes := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
		names := make([]string, 0, 128)
		for note := 0; note < 128; note++ {
			names = append(names, fmt.Sprintf("%s%d", prefixes[note%12], (note/12)-1))
		}
		r.AddData(identifier+"_names", names, config.DataAttributes, nil)
	}

	if slt.IsSelected("octaves") {
		octaves := make([]uint8, 0, 128)
		for note := 0; note < 128; note++ {
			octaves = append(octaves, uint8(note/12))
		}
		r.AddData(identifier+"_octaves", octaves, config.DataAttributes, nil)
	}

	return nil
}
