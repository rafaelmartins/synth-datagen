package filters

import (
	"fmt"
	"math"

	"rafaelmartins.com/p/synth-datagen/internal/convert"
	"rafaelmartins.com/p/synth-datagen/internal/datareg"
	"rafaelmartins.com/p/synth-datagen/internal/renderer"
	"rafaelmartins.com/p/synth-datagen/internal/selector"
)

type Filters struct{}

type filtersConfig struct {
	SampleRate                            float64
	DataAttributes                        []string
	Frequencies                           int
	FrequencyMax                          float64
	FrequencyMin                          float64
	FrequencyDescriptionsStringWidth      *int
	CoefficientsOnepoleScalarType         *string `selectors:"lowpass_onepole,highpass_onepole"`
	CoefficientsOnepoleFractionalBitWidth *uint8
}

func (*Filters) GetName() string {
	return "filters"
}

func (*Filters) GetAllowedSelectors() []string {
	return []string{"lowpass_onepole", "highpass_onepole", "descriptions"}
}

type filter1Pole struct {
	A1 float64
	B0 float64
	B1 float64
}

func (f *Filters) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]any, slt *selector.Selector) error {
	config := filtersConfig{}
	if err := dreg.Evaluate(f.GetName(), &config, pmt, slt); err != nil {
		return err
	}

	tmpFreqs := make([]float64, 0, config.Frequencies)
	for i := 0.; i < float64(config.Frequencies); i++ {
		tmpFreqs = append(tmpFreqs, -1+math.Exp(3*i/float64(config.Frequencies-1)))
	}

	dF := config.FrequencyMax - config.FrequencyMin
	freqs := make([]float64, 0, config.Frequencies)
	nfreqs := make([]float64, 0, config.Frequencies)
	for _, fr := range tmpFreqs {
		freq := config.FrequencyMin + dF*fr/tmpFreqs[config.Frequencies-1]
		freqs = append(freqs, freq)
		nfreqs = append(nfreqs, freq/config.SampleRate)
	}

	bw := 0
	if config.CoefficientsOnepoleFractionalBitWidth != nil {
		bw = int(*config.CoefficientsOnepoleFractionalBitWidth)
	}

	if slt.IsSelected("lowpass_onepole") {
		lp := make([]filter1Pole, 0, config.Frequencies)
		for _, freq := range nfreqs {
			a1 := (1. - math.Tan(math.Pi*freq)) / (1. + math.Tan(math.Pi*freq))
			b0 := (1 - a1) / 2
			lp = append(lp, filter1Pole{
				A1: a1 * float64(int(1)<<bw),
				B0: b0 * float64(int(1)<<bw),
				B1: b0 * float64(int(1)<<bw),
			})
		}
		v, err := convert.SliceStruct(lp, *config.CoefficientsOnepoleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_lowpass_onepole_coefficients", v, config.DataAttributes, nil)
	}

	if slt.IsSelected("highpass_onepole") {
		hp := make([]filter1Pole, 0, config.Frequencies)
		for _, freq := range nfreqs {
			a1 := (1. - math.Tan(math.Pi*freq)) / (1. + math.Tan(math.Pi*freq))
			b0 := (1 + a1) / 2
			hp = append(hp, filter1Pole{
				A1: a1 * float64(int(1)<<bw),
				B0: b0 * float64(int(1)<<bw),
				B1: -b0 * float64(int(1)<<bw),
			})
		}
		v, err := convert.SliceStruct(hp, *config.CoefficientsOnepoleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_highpass_onepole_coefficients", v, config.DataAttributes, nil)
	}

	if slt.IsSelected("descriptions") {
		desc := make([]string, 0, config.Frequencies)
		for _, freq := range freqs {
			if freq > 1000 {
				desc = append(desc, fmt.Sprintf("%.2fkHz", freq/1000))
			} else {
				desc = append(desc, fmt.Sprintf("%dHz", int(freq)))
			}
		}
		r.AddData(identifier+"_frequency_descriptions", desc, config.DataAttributes, config.FrequencyDescriptionsStringWidth)
	}

	return nil
}
