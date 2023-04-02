package modules

import (
	"errors"
	"fmt"

	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/modules/blwavetables"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/selector"
)

type Module interface {
	GetName() string
	GetAllowedSelectors() []string
	Render(r renderer.Renderer, identifier string, dreg *datareg.DataReg, pmt map[string]interface{}, slt *selector.Selector) error
}

var (
	mreg = []Module{
		&blwavetables.BandLimitedWavetables{},
	}

	dreg = &datareg.DataReg{}
)

func SetGlobalParameters(pmt map[string]interface{}) {
	dreg = datareg.New(pmt)
}

func Render(r renderer.Renderer, identifier string, module string, pmt map[string]interface{}, sel []string) error {
	if r == nil {
		return errors.New("modules: header not defined")
	}

	for _, mod := range mreg {
		if mod != nil && mod.GetName() == module {
			slt, err := selector.New(mod.GetAllowedSelectors(), sel)
			if err != nil {
				return err
			}
			return mod.Render(r, identifier, dreg, pmt, slt)
		}
	}
	return fmt.Errorf("modules: module not found: %s", module)
}
