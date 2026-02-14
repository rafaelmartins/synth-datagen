package adsr

import (
	"errors"
	"fmt"
	"math"

	"rafaelmartins.com/p/synth-datagen/internal/convert"
	"rafaelmartins.com/p/synth-datagen/internal/datareg"
	"rafaelmartins.com/p/synth-datagen/internal/renderer"
	"rafaelmartins.com/p/synth-datagen/internal/selector"
)

const (
	as3310AttackAsymptoteVoltage = 7.0
	as3310AttackPeakVoltage      = 5.0
)

type ADSR struct{}

type adsrConfig struct {
	Samples                      int
	DataAttributes               []string
	SampleAmplitude              *float64 `selectors:"curves_as3310,curves_linear"`
	SampleScalarType             *string  `selectors:"curves_as3310,curves_linear"`
	SampleRate                   *float64 `selectors:"time_steps"`
	TimeSteps                    *int     `selectors:"time_steps,descriptions"`
	TimeStepsMinMs               *int     `selectors:"time_steps,descriptions"`
	TimeStepsMaxMs               *int     `selectors:"time_steps,descriptions"`
	TimeStepsScalarType          *string  `selectors:"time_steps"`
	TimeStepsFractionalBitWidth  *uint8
	LevelDescriptions            *int `selectors:"descriptions"`
	LevelDescriptionsStringWidth *int
	TimeDescriptionsStringWidth  *int
}

func (*ADSR) GetName() string {
	return "adsr"
}

func (*ADSR) GetAllowedSelectors() []string {
	return []string{"curves_as3310", "curves_linear", "time_steps", "descriptions"}
}

func (a *ADSR) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	config := adsrConfig{}
	if err := dreg.Evaluate(a.GetName(), &config, pmt, slt); err != nil {
		return err
	}

	sampleBase := make([]float64, 0, config.Samples)
	for i := 0; i < config.Samples; i++ {
		sampleBase = append(sampleBase, float64(i)/(float64(config.Samples-1)))
	}

	if slt.IsSelected("curves_as3310") {
		baseCurve := make([]float64, 0, config.Samples)
		for _, t := range sampleBase {
			baseCurve = append(baseCurve, 1.-math.Exp(-3*t))
		}

		target := baseCurve[config.Samples-1]
		if target == 0 {
			return errors.New("modules: adsr: base curve target is zero")
		}

		attackPeak := 0.
		for i, v := range baseCurve {
			if v/target >= as3310AttackPeakVoltage/as3310AttackAsymptoteVoltage {
				attackPeak = sampleBase[i]
				break
			}
		}

		baseAttackCurve := make([]float64, 0, config.Samples)
		for _, t := range sampleBase {
			baseAttackCurve = append(baseAttackCurve, 1.-math.Exp(-3*t*attackPeak))
		}

		attackCurve := make([]float64, 0, config.Samples)
		releaseCurve := make([]float64, 0, config.Samples)
		for i := 0; i < config.Samples; i++ {
			attackCurve = append(attackCurve, *config.SampleAmplitude*baseAttackCurve[i]/baseAttackCurve[config.Samples-1])
			releaseCurve = append(releaseCurve, *config.SampleAmplitude*baseCurve[i]/target)
		}

		atk, err := convert.Slice(attackCurve, *config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_as3310_attack", atk, config.DataAttributes, nil)

		rel, err := convert.Slice(releaseCurve, *config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_as3310_decay_release", rel, config.DataAttributes, nil)
	}

	if slt.IsSelected("curves_linear") {
		linearCurve := make([]float64, 0, config.Samples)
		for _, t := range sampleBase {
			linearCurve = append(linearCurve, *config.SampleAmplitude*t)
		}

		lin, err := convert.Slice(linearCurve, *config.SampleScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_linear", lin, config.DataAttributes, nil)
	}

	times := []float64{}

	if slt.IsSelected("time_steps") || slt.IsSelected("descriptions") {
		tmpTimes := make([]float64, 0, *config.TimeSteps)
		for i := 0.; i < float64(*config.TimeSteps); i++ {
			tmpTimes = append(tmpTimes, -1.+math.Exp(6*i/(float64(*config.TimeSteps-1))))
		}

		dT := float64(*config.TimeStepsMaxMs - *config.TimeStepsMinMs)
		times = make([]float64, 0, *config.TimeSteps)
		for _, t := range tmpTimes {
			times = append(times, float64(*config.TimeStepsMinMs)+(dT*t/tmpTimes[*config.TimeSteps-1]))
		}
	}

	if slt.IsSelected("time_steps") {
		timeSteps := make([]float64, 0, *config.TimeSteps)
		for _, t := range times {
			timeSteps = append(timeSteps, (float64(config.Samples)*1000.)/(t*(*config.SampleRate)))
		}

		if config.TimeStepsFractionalBitWidth != nil {
			for idx := range timeSteps {
				timeSteps[idx] *= float64(int(1) << *config.TimeStepsFractionalBitWidth)
			}
		}

		ts, err := convert.Slice(timeSteps, *config.TimeStepsScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_time_steps", ts, config.DataAttributes, nil)
	}

	if slt.IsSelected("descriptions") {
		levels := make([]string, 0, *config.LevelDescriptions)
		for i := 0.; i < float64(*config.LevelDescriptions); i++ {
			levels = append(levels, fmt.Sprintf("%.1f%%", 100.*i/float64(*config.LevelDescriptions-1)))
		}
		r.AddData(identifier+"_level_descriptions", levels, config.DataAttributes, config.LevelDescriptionsStringWidth)

		timed := make([]string, 0, *config.TimeSteps)
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
		r.AddData(identifier+"_time_descriptions", timed, config.DataAttributes, config.TimeDescriptionsStringWidth)
	}

	return nil
}
