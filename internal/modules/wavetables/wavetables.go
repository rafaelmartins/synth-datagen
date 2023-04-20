package wavetables

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

type Wavetables struct{}

type wavetablesConfig struct {
	SamplesPerCycle            int
	SampleAmplitude            float64
	SampleScalarType           string
	DataAttributes             []string
	SampleRate                 *float64 `selectors:"blsquare,bltriangle,blsawtooth"`
	BandlimitedOmitHighOctaves *int
}

func (*Wavetables) GetName() string {
	return "wavetables"
}

func (*Wavetables) GetAllowedSelectors() []string {
	return []string{"sine", "blsquare", "bltriangle", "blsawtooth"}
}

func (bl *Wavetables) fixWavetable(config *wavetablesConfig, data []float64) []float64 {
	min := 0.
	max := 0.
	if len(data) > 0 {
		min = data[0]
		max = data[0]
	}

	for i := 1; i < len(data); i++ {
		if data[i] < min {
			min = data[i]
		}
		if data[i] > max {
			max = data[i]
		}
	}

	scaleFactor := (2 * config.SampleAmplitude) / math.Abs(max-min)
	rv := make([]float64, config.SamplesPerCycle)
	for i := range data {
		rv[len(data)-i-1] = (data[i]-min)*scaleFactor - config.SampleAmplitude
	}
	return rv
}

func (bl *Wavetables) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	config := wavetablesConfig{}
	if err := dreg.Evaluate(bl.GetName(), &config, pmt, slt); err != nil {
		return err
	}

	amp, err := convert.Scalar(config.SampleAmplitude, config.SampleScalarType)
	if err != nil {
		return err
	}
	r.AddMacro(identifier+"_sample_amplitude", amp, true, false)

	if slt.IsSelected("sine") {
		sine := make([]float64, 0, config.SamplesPerCycle)
		for i := 0; i < config.SamplesPerCycle; i++ {
			sine = append(sine, float64(config.SampleAmplitude*math.Sin(2*math.Pi*float64(i)/float64(config.SamplesPerCycle))))
		}
		v, err := convert.Slice(sine, config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_sine", v, config.DataAttributes, nil)
	}

	if slt.IsSelected("blsquare") || slt.IsSelected("bltriangle") || slt.IsSelected("blsawtooth") {
		numOctaves := int(math.Ceil(128.0 / 12))
		if config.BandlimitedOmitHighOctaves != nil {
			if *config.BandlimitedOmitHighOctaves < 0 || *config.BandlimitedOmitHighOctaves >= numOctaves {
				return fmt.Errorf("wavetables: bandlimited_omit_high_octaves must be >= 0 and < %d", numOctaves)
			}
			numOctaves -= *config.BandlimitedOmitHighOctaves
		}

		squares := make([][]float64, 0, numOctaves)
		triangles := make([][]float64, 0, numOctaves)
		sawtooths := make([][]float64, 0, numOctaves)

		for oct := 0; oct < numOctaves; oct++ {
			freq := wavetableFrequency(oct)
			period := *config.SampleRate / freq
			harmonics := float64(int(period))
			if math.Mod(harmonics, 2) == 0 {
				harmonics--
			}

			blit := make([]float64, 0, config.SamplesPerCycle)
			for i := 0; i < config.SamplesPerCycle; i++ {
				normalizedPos := (float64(i) - float64(config.SamplesPerCycle)/2) / float64(config.SamplesPerCycle)
				if normalizedPos == 0 {
					blit = append(blit, 1.0)
				} else {
					blit = append(blit, math.Sin(math.Pi*normalizedPos*harmonics)/(harmonics*math.Sin(math.Pi*normalizedPos)))
				}
			}
			blitMid := config.SamplesPerCycle / 2

			if slt.IsSelected("blsquare") || slt.IsSelected("bltriangle") {
				square := make([]float64, 0, config.SamplesPerCycle)
				squareSum := 0.
				v := 0.
				for i := 0; i < config.SamplesPerCycle; i++ {
					v += blit[i]
					if i < blitMid {
						v -= blit[i+blitMid]
					} else {
						v -= blit[i-blitMid]
					}
					square = append(square, v)
					squareSum += v
				}
				squareAvg := squareSum / float64(config.SamplesPerCycle)
				squares = append(squares, bl.fixWavetable(&config, square))

				if slt.IsSelected("bltriangle") {
					triangle := make([]float64, 0, config.SamplesPerCycle)
					v = 0.
					for _, sq := range square {
						v += sq - squareAvg
						triangle = append(triangle, v)
					}
					triangleStart := config.SamplesPerCycle / 4
					triangle = append(triangle[triangleStart:], triangle[:triangleStart]...)
					triangles = append(triangles, bl.fixWavetable(&config, triangle))
				}
			}

			if slt.IsSelected("blsawtooth") {
				sawtooth := make([]float64, 0, config.SamplesPerCycle)
				v := 0.
				for i := 0; i < config.SamplesPerCycle; i++ {
					v -= 1. / period
					if i < blitMid {
						v += blit[i+blitMid]
					} else {
						v += blit[i-blitMid]
					}
					sawtooth = append(sawtooth, -v)
				}
				sawtooths = append(sawtooths, bl.fixWavetable(&config, sawtooth))
			}
		}

		if slt.IsSelected("blsquare") {
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range squares {
				v, err := convert.Slice(wt, config.SampleScalarType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_blsquare", rv, config.DataAttributes, nil)
		}

		if slt.IsSelected("bltriangle") {
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range triangles {
				v, err := convert.Slice(wt, config.SampleScalarType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_bltriangle", rv, config.DataAttributes, nil)
		}

		if slt.IsSelected("blsawtooth") {
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range sawtooths {
				v, err := convert.Slice(wt, config.SampleScalarType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_blsawtooth", rv, config.DataAttributes, nil)
		}
	}

	return nil
}
