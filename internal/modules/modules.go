package modules

import (
	"errors"
	"fmt"

	"github.com/rafaelmartins/synth-datagen/internal/datareg"
	"github.com/rafaelmartins/synth-datagen/internal/modules/blwavetables"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
)

type Module interface {
	GetName() string
	Render(r renderer.Renderer, dreg *datareg.DataReg, pmt map[string]interface{}) error
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

func Render(r renderer.Renderer, module string, pmt map[string]interface{}) error {
	if r == nil {
		return errors.New("modules: header not defined")
	}

	for _, mod := range mreg {
		if mod != nil && mod.GetName() == module {
			return mod.Render(r, dreg, pmt)
		}
	}
	return fmt.Errorf("modules: module not found: %s", module)
}
