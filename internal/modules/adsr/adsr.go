package adsr

import (
	"fmt"
	"math"

	"github.com/rafaelmartins/synth-datagen/internal/convert"
	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
	"github.com/rafaelmartins/synth-datagen/internal/utils"
)

const (
	as3310AttackAsymptoteVoltage = 7.0
	as3310AttackPeakVoltage      = 5.0
)

type ADSR struct {
	config struct {
		CurveSamples          int
		TimeSamples           int
		TimeMinMs             int
		TimeMaxMs             int
		DataAttributes        []string
		DataScalarType        *string  `selectors:"curves_as3310,curves_linear"`
		SampleAmplitude       *float64 `selectors:"curves_as3310,curves_linear"`
		SampleRate            *float64 `selectors:"steps"`
		SamplesPerCycle       *int     `selectors:"steps"`
		LevelSamples          *int     `selectors:"descriptions"`
		LevelDescriptionWidth *int
		TimeDescriptionWidth  *int
	}
}

func (*ADSR) GetName() string {
	return "adsr"
}

func (*ADSR) GetAllowedSelectors() []string {
	return []string{"curves_as3310", "curves_linear", "steps", "descriptions"}
}

func (a *ADSR) Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error {
	if err := dreg.Evaluate(&a.config, pmt, slt); err != nil {
		return err
	}

	sampleBase := make([]float64, 0, a.config.CurveSamples)
	for i := 0.; i < float64(a.config.CurveSamples); i++ {
		sampleBase = append(sampleBase, i/(float64(a.config.CurveSamples-1)))
	}

	if slt.IsSelected("curves_as3310") {
		baseCurve := make([]float64, 0, a.config.CurveSamples)
		for _, t := range sampleBase {
			baseCurve = append(baseCurve, 1.-math.Exp(-3*t))
		}

		attackPeak := 0.
		for i, v := range baseCurve {
			if v/baseCurve[a.config.CurveSamples-1] >= as3310AttackPeakVoltage/as3310AttackAsymptoteVoltage {
				attackPeak = sampleBase[i]
				break
			}
		}

		baseAttackCurve := make([]float64, 0, a.config.CurveSamples)
		for _, t := range sampleBase {
			baseAttackCurve = append(baseAttackCurve, 1.-math.Exp(-3*t*attackPeak))
		}

		attackCurve := make([]float64, 0, a.config.CurveSamples)
		releaseCurve := make([]float64, 0, a.config.CurveSamples)
		for i := 0; i < a.config.CurveSamples; i++ {
			attackCurve = append(attackCurve, *a.config.SampleAmplitude*baseAttackCurve[i]/baseAttackCurve[a.config.CurveSamples-1])
			releaseCurve = append(releaseCurve, *a.config.SampleAmplitude*baseCurve[i]/baseCurve[a.config.CurveSamples-1])
		}

		atk, err := convert.Slice(attackCurve, *a.config.DataScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_attack", atk, a.config.DataAttributes)

		rel, err := convert.Slice(releaseCurve, *a.config.DataScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_decay_release", rel, a.config.DataAttributes)
	}

	if slt.IsSelected("curves_linear") {
		linearCurve := make([]float64, 0, a.config.CurveSamples)
		for _, t := range sampleBase {
			linearCurve = append(linearCurve, *a.config.SampleAmplitude*t)
		}

		lin, err := convert.Slice(linearCurve, *a.config.DataScalarType)
		if err != nil {
			return err
		}
		r.AddData(identifier+"_curve_linear", lin, a.config.DataAttributes)
	}

	tmpTimes := make([]float64, 0, a.config.TimeSamples)
	for i := 0.; i < float64(a.config.TimeSamples); i++ {
		tmpTimes = append(tmpTimes, -1.+math.Exp(6*i/(float64(a.config.TimeSamples-1))))
	}

	dT := float64(a.config.TimeMaxMs - a.config.TimeMinMs)
	times := make([]float64, 0, a.config.TimeSamples)
	for _, t := range tmpTimes {
		times = append(times, float64(a.config.TimeMinMs)+(dT*t/tmpTimes[a.config.TimeSamples-1]))
	}

	if slt.IsSelected("steps") {
		timeSteps := make([]uint32, 0, a.config.TimeSamples)
		for _, t := range times {
			timeSteps = append(timeSteps, uint32(((float64(*a.config.SamplesPerCycle)*1000.)/(t*(*a.config.SampleRate)))*(1<<16)))
		}
		r.AddData(identifier+"_time_steps", timeSteps, a.config.DataAttributes)
	}

	if slt.IsSelected("descriptions") {
		levels := make([]string, 0, *a.config.LevelSamples)
		for i := 0.; i < float64(*a.config.LevelSamples); i++ {
			levels = append(levels, fmt.Sprintf("%.1f%%", 100.*i/float64(*a.config.LevelSamples-1)))
		}
		if a.config.LevelDescriptionWidth != nil {
			l, err := utils.FormatStringSliceWidth(levels, *a.config.LevelDescriptionWidth)
			if err != nil {
				return fmt.Errorf("adsr: level_descriptions: %w", err)
			}
			levels = l
			r.AddMacro(identifier+"_level_descriptions_width", int(math.Abs(float64(*a.config.LevelDescriptionWidth))), false, false)
		}
		r.AddData(identifier+"_level_descriptions", levels, a.config.DataAttributes)

		timed := make([]string, 0, a.config.TimeSamples)
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
		if a.config.TimeDescriptionWidth != nil {
			l, err := utils.FormatStringSliceWidth(timed, *a.config.TimeDescriptionWidth)
			if err != nil {
				return fmt.Errorf("adsr: time_descriptions: %w", err)
			}
			timed = l
			r.AddMacro(identifier+"_time_descriptions_width", int(math.Abs(float64(*a.config.TimeDescriptionWidth))), false, false)
		}
		r.AddData(identifier+"_time_descriptions", timed, a.config.DataAttributes)
	}

	return nil
}
