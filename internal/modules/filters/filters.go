package filters

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

type Filters struct{}

type filtersConfig struct {
	SampleRate                       float64
	DataAttributes                   []string
	Frequencies                      int
	FrequencyMax                     float64
	FrequencyMin                     float64
	FrequencyDescriptionsStringWidth *int
}

func (*Filters) GetName() string {
	return "filters"
}

func (*Filters) GetAllowedSelectors() []string {
	return []string{"lowpass_1pole", "highpass_1pole", "descriptions"}
}

type filter1Pole struct {
	A1 int8
	B0 int8
	B1 int8
}

func (f *Filters) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
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
	alphas := make([]float64, 0, config.Frequencies)
	for _, fr := range tmpFreqs {
		frq := config.FrequencyMin + dF*fr/tmpFreqs[config.Frequencies-1]
		freqs = append(freqs, frq)
		alphas = append(alphas, 2*math.Pi*frq/config.SampleRate)
	}

	if slt.IsSelected("lowpass_1pole") {
		lp := make([]filter1Pole, 0, config.Frequencies)
		for _, alpha := range alphas {
			lp = append(lp, filter1Pole{
				A1: int8((-(alpha - 2) / (alpha + 2)) * (1 << 7)),
				B0: int8((alpha / (alpha + 2)) * (1 << 7)),
				B1: int8((alpha / (alpha + 2)) * (1 << 7)),
			})
		}
		r.AddData(identifier+"_lowpass_1pole_coefficients", lp, config.DataAttributes, nil)
	}

	if slt.IsSelected("highpass_1pole") {
		hp := make([]filter1Pole, 0, config.Frequencies)
		for _, alpha := range alphas {
			hp = append(hp, filter1Pole{
				A1: int8(((1 - alpha/2) / (1 + alpha/2)) * (1 << 7)),
				B0: int8((1 / (1 + alpha/2)) * (1 << 7)),
				B1: int8((-1 / (1 + alpha/2)) * (1 << 7)),
			})
		}
		r.AddData(identifier+"_highpass_1pole_coefficients", hp, config.DataAttributes, nil)
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
