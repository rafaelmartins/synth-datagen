package codegen

import (
	"fmt"
	"io"
	"strings"

	"rafaelmartins.com/p/synth-datagen/internal/codegen/stringify"
	"rafaelmartins.com/p/synth-datagen/internal/utils"
)

type data struct {
	identifier string
	value      interface{}
	attributes []string
	strWidth   *int
}

type dataList []*data

func (d *dataList) add(identifier string, value interface{}, attributes []string, strWidth *int) {
	*d = append(*d, &data{
		identifier: identifier,
		value:      value,
		attributes: attributes,
		strWidth:   strWidth,
	})
}

func (d dataList) write(w io.Writer) error {
	for _, dat := range d {
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}

		if dat.strWidth != nil {
			if _, ok := dat.value.(string); ok {
				if err := utils.ApplyStringWidth(&dat.value, *dat.strWidth); err != nil {
					return err
				}
			} else {
				if err := utils.ApplyStringWidth(dat.value, *dat.strWidth); err != nil {
					return err
				}
			}
		}

		value, ctype, dim, err := stringify.Stringify(dat.value)
		if err != nil {
			return err
		}

		if ctype == "char*" && dat.strWidth != nil {
			ctype = "char"
			dim = append(dim, utils.Abs(*dat.strWidth))
		}

		ctyped := "static const " + ctype + " " + dat.identifier
		for _, d := range dim {
			ctyped += fmt.Sprintf("[%d]", d)
		}
		ctyped += " "
		if len(dat.attributes) > 0 {
			ctyped += strings.Join(dat.attributes, " ") + " "
		}
		ctyped += "= " + value + ";\n"

		if _, err := io.WriteString(w, ctyped); err != nil {
			return err
		}

		switch len(dim) {
		case 0:
		case 1:
			if _, err := fmt.Fprintf(w, "#define %s_len %d\n", dat.identifier, dim[0]); err != nil {
				return err
			}

		case 2:
			if _, err := fmt.Fprintf(w, "#define %s_rows %d\n", dat.identifier, dim[0]); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "#define %s_cols %d\n", dat.identifier, dim[1]); err != nil {
				return err
			}

		default:
			for i, d := range dim {
				if _, err := fmt.Fprintf(w, "#define %s_len_%d %d\n", dat.identifier, i, d); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
