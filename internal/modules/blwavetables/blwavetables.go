package blwavetables

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
)

type BandLimitedWavetables struct {
	config struct {
		ClockFrequency          float64
		SampleRate              float64
		WaveformAmplitude       float64
		WaveformSamplesPerCycle int
		Waveforms               []string
		WavetableDataType       string
		WavetableAttributes     []string
	}
}

func (*BandLimitedWavetables) GetName() string {
	return "blwavetables"
}

func (bl *BandLimitedWavetables) fixWavetable(data []float64) []float64 {
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

	scaleFactor := (2 * bl.config.WaveformAmplitude) / math.Abs(max-min)
	rv := make([]float64, bl.config.WaveformSamplesPerCycle)
	for i := range data {
		rv[len(data)-i-1] = (data[i]-min)*scaleFactor - bl.config.WaveformAmplitude
	}
	return rv
}

func (bl *BandLimitedWavetables) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}) error {
	if err := dreg.Evaluate(&bl.config, pmt); err != nil {
		return err
	}

	sine := make([]float64, 0, bl.config.WaveformSamplesPerCycle)
	for i := 0; i < bl.config.WaveformSamplesPerCycle; i++ {
		sine = append(sine, float64(bl.config.WaveformAmplitude*math.Sin(2*math.Pi*float64(i)/float64(bl.config.WaveformSamplesPerCycle))))
	}

	numOctaves := int(math.Ceil(128.0 / 12))

	squares := make([][]float64, 0, numOctaves)
	triangles := make([][]float64, 0, numOctaves)
	sawtooths := make([][]float64, 0, numOctaves)

	for oct := 0; oct < numOctaves-1; oct++ {
		freq := wavetableFrequency(oct)
		period := bl.config.SampleRate / freq
		harmonics := float64(int(period))
		if math.Mod(harmonics, 2) == 0 {
			harmonics--
		}

		blit := make([]float64, 0, bl.config.WaveformSamplesPerCycle)
		for i := 0; i < bl.config.WaveformSamplesPerCycle; i++ {
			normalizedPos := (float64(i) - float64(bl.config.WaveformSamplesPerCycle)/2) / float64(bl.config.WaveformSamplesPerCycle)
			if normalizedPos == 0 {
				blit = append(blit, 1.0)
			} else {
				blit = append(blit, math.Sin(math.Pi*normalizedPos*harmonics)/(harmonics*math.Sin(math.Pi*normalizedPos)))
			}
		}
		blitMid := bl.config.WaveformSamplesPerCycle / 2

		square := make([]float64, 0, bl.config.WaveformSamplesPerCycle)
		squareSum := 0.
		v := 0.
		for i := 0; i < bl.config.WaveformSamplesPerCycle; i++ {
			v += blit[i]
			if i < blitMid {
				v -= blit[i+blitMid]
			} else {
				v -= blit[i-blitMid]
			}
			square = append(square, v)
			squareSum += v
		}
		squareAvg := squareSum / float64(bl.config.WaveformSamplesPerCycle)
		squares = append(squares, bl.fixWavetable(square))

		triangle := make([]float64, 0, bl.config.WaveformSamplesPerCycle)
		v = 0.
		for _, sq := range square {
			v += sq - squareAvg
			triangle = append(triangle, v)
		}
		triangleStart := bl.config.WaveformSamplesPerCycle / 4
		triangle = append(triangle[triangleStart:], triangle[:triangleStart]...)
		triangles = append(triangles, bl.fixWavetable(triangle))

		sawtooth := make([]float64, 0, bl.config.WaveformSamplesPerCycle)
		v = 0.
		for i := 0; i < bl.config.WaveformSamplesPerCycle; i++ {
			v -= 1. / period
			if i < blitMid {
				v += blit[i+blitMid]
			} else {
				v += blit[i-blitMid]
			}
			sawtooth = append(sawtooth, -v)
		}
		sawtooths = append(sawtooths, bl.fixWavetable(sawtooth))
	}

	squares = append(squares, sine)
	triangles = append(triangles, sine)
	sawtooths = append(sawtooths, sine)

	for _, wf := range bl.config.Waveforms {
		switch wf {
		case "sine":
			v, err := convert.Slice(sine, bl.config.WavetableDataType)
			if err != nil {
				return err
			}
			r.AddData(identifier+"_sine", v, bl.config.WavetableAttributes)

		case "square":
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range squares {
				v, err := convert.Slice(wt, bl.config.WavetableDataType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_square", rv, bl.config.WavetableAttributes)

		case "triangle":
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range triangles {
				v, err := convert.Slice(wt, bl.config.WavetableDataType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_triangle", rv, bl.config.WavetableAttributes)

		case "sawtooth":
			rv := make([]interface{}, 0, numOctaves)
			for _, wt := range sawtooths {
				v, err := convert.Slice(wt, bl.config.WavetableDataType)
				if err != nil {
					return err
				}
				rv = append(rv, v)
			}
			r.AddData(identifier+"_sawtooth", rv, bl.config.WavetableAttributes)

		default:
			return fmt.Errorf("blwavetables: invalid waveform: %s", wf)
		}
	}

	return nil
}
