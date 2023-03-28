package codegen

import (
	"fmt"
	"io"

	"github.com/rafaelmartins/synth-datagen/internal/codegen/stringify"
)

type macro struct {
	identifier string
	value      interface{}
	hex        bool
}

type macroList []*macro

func (m *macroList) add(identifier string, value interface{}, hex bool) {
	*m = append(*m, &macro{
		identifier: identifier,
		value:      value,
		hex:        hex,
	})
}

func (m macroList) write(w io.Writer) error {
	if len(m) > 0 {
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	for _, mac := range m {
		val, err := stringify.StringifyValue(mac.value, mac.hex)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "#define %s %s\n", mac.identifier, val); err != nil {
			return err
		}
	}

	return nil
}
