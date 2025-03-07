package codegen

import (
	"fmt"
	"io"

	"rafaelmartins.com/p/synth-datagen/internal/codegen/stringify"
)

type macro struct {
	identifier string
	value      interface{}
	hex        bool
	raw        bool
}

type macroList []*macro

func (m *macroList) add(identifier string, value interface{}, hex bool, raw bool) {
	*m = append(*m, &macro{
		identifier: identifier,
		value:      value,
		hex:        hex,
		raw:        raw,
	})
}

func (m macroList) write(w io.Writer) error {
	if len(m) > 0 {
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	for _, mac := range m {
		if mac.raw {
			if _, err := fmt.Fprintf(w, "#define %s %v\n", mac.identifier, mac.value); err != nil {
				return err
			}
			continue
		}

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
