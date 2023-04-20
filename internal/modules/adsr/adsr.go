package adsr

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

const (
	as3310AttackAsymptoteVoltage = 7.0
	as3310AttackPeakVoltage      = 5.0
)

type ADSR struct {
	config struct {
		Samples                      int
		DataAttributes               []string
		SampleAmplitude              *float64 `selectors:"curves_as3310,curves_linear"`
		SampleScalarType             *string  `selectors:"curves_as3310,curves_linear"`
		SampleRate                   *float64 `selectors:"time_steps"`
		SamplesPerCycle              *int     `selectors:"time_steps"`
		TimeSteps                    *int     `selectors:"time_steps,descriptions"`
		TimeStepsMinMs               *int     `selectors:"time_steps,descriptions"`
		TimeStepsMaxMs               *int     `selectors:"time_steps,descriptions"`
		TimeStepsScalarType          *string  `selectors:"time_steps"`
		TimeStepsFractionalBitWidth  *uint8   `selectors:"time_steps"`
		LevelDescriptions            *int     `selectors:"descriptions"`
		LevelDescriptionsStringWidth *int
		TimeDescriptionsStringWidth  *int
	}
}

func (*ADSR) GetName() string {
	return "adsr"
}

func (*ADSR) GetAllowedSelectors() []string {
	return []string{"curves_as3310", "curves_linear", "time_steps", "descriptions"}
}

func (a *ADSR) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	if err := dreg.Evaluate(a.GetName(), &a.config, pmt, slt); err != nil {
		return err
	}

	sampleBase := make([]float64, 0, a.config.Samples)
	for i := 0.; i < float64(a.config.Samples); i++ {
		sampleBase = append(sampleBase, i/(float64(a.config.Samples-1)))
	}

	if slt.IsSelected("curves_as3310") {
		baseCurve := make([]float64, 0, a.config.Samples)
		for _, t := range sampleBase {
			baseCurve = append(baseCurve, 1.-math.Exp(-3*t))
		}

		attackPeak := 0.
		for i, v := range baseCurve {
			if v/baseCurve[a.config.Samples-1] >= as3310AttackPeakVoltage/as3310AttackAsymptoteVoltage {
				attackPeak = sampleBase[i]
				break
			}
		}

		baseAttackCurve := make([]float64, 0, a.config.Samples)
		for _, t := range sampleBase {
			baseAttackCurve = append(baseAttackCurve, 1.-math.Exp(-3*t*attackPeak))
		}

		attackCurve := make([]float64, 0, a.config.Samples)
		releaseCurve := make([]float64, 0, a.config.Samples)
		for i := 0; i < a.config.Samples; i++ {
			attackCurve = append(attackCurve, *a.config.SampleAmplitude*baseAttackCurve[i]/baseAttackCurve[a.config.Samples-1])
			releaseCurve = append(releaseCurve, *a.config.SampleAmplitude*baseCurve[i]/baseCurve[a.config.Samples-1])
		}

		atk, err := convert.Slice(attackCurve, *a.config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_attack", atk, a.config.DataAttributes, nil)

		rel, err := convert.Slice(releaseCurve, *a.config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_decay_release", rel, a.config.DataAttributes, nil)
	}

	if slt.IsSelected("curves_linear") {
		linearCurve := make([]float64, 0, a.config.Samples)
		for _, t := range sampleBase {
			linearCurve = append(linearCurve, *a.config.SampleAmplitude*t)
		}

		lin, err := convert.Slice(linearCurve, *a.config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_linear", lin, a.config.DataAttributes, nil)
	}

	if slt.IsSelected("curves_as3310") || slt.IsSelected("curves_linear") {
		amp, err := convert.Scalar(*a.config.SampleAmplitude, *a.config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddMacro(identifier+"_sample_amplitude", amp, true, false)
	}

	times := []float64{}

	if slt.IsSelected("time_steps") || slt.IsSelected("descriptions") {
		tmpTimes := make([]float64, 0, *a.config.TimeSteps)
		for i := 0.; i < float64(*a.config.TimeSteps); i++ {
			tmpTimes = append(tmpTimes, -1.+math.Exp(6*i/(float64(*a.config.TimeSteps-1))))
		}

		dT := float64(*a.config.TimeStepsMaxMs - *a.config.TimeStepsMinMs)
		times = make([]float64, 0, *a.config.TimeSteps)
		for _, t := range tmpTimes {
			times = append(times, float64(*a.config.TimeStepsMinMs)+(dT*t/tmpTimes[*a.config.TimeSteps-1]))
		}
	}

	if slt.IsSelected("time_steps") {
		timeSteps := make([]uint32, 0, *a.config.TimeSteps)
		for _, t := range times {
			timeSteps = append(timeSteps, uint32(((float64(*a.config.SamplesPerCycle)*1000.)/(t*(*a.config.SampleRate)))*float64(int(1)<<*a.config.TimeStepsFractionalBitWidth)))
		}

		ts, err := convert.Slice(timeSteps, *a.config.TimeStepsScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_time_steps", ts, a.config.DataAttributes, nil)
	}

	if slt.IsSelected("descriptions") {
		levels := make([]string, 0, *a.config.LevelDescriptions)
		for i := 0.; i < float64(*a.config.LevelDescriptions); i++ {
			levels = append(levels, fmt.Sprintf("%.1f%%", 100.*i/float64(*a.config.LevelDescriptions-1)))
		}
		r.AddData(identifier+"_level_descriptions", levels, a.config.DataAttributes, a.config.LevelDescriptionsStringWidth)

		timed := make([]string, 0, *a.config.TimeSteps)
		for _, t := range times {
			if t > 10000 {
				timed = append(timed, fmt.Sprintf("%.1fs", t/1000))
				continue
			}
			if t > 1000 {
				timed = append(timed, fmt.Sprintf("%.2fs", t/1000))
				continue
			}
			timed = append(timed, fmt.Sprintf("%dms", int(t)))
		}
		r.AddData(identifier+"_time_descriptions", timed, a.config.DataAttributes, a.config.TimeDescriptionsStringWidth)
	}

	return nil
}
