package stringify

import (
	"errors"
	"reflect"
	"strings"

	"rafaelmartins.com/p/synth-datagen/internal/ctypes"
)

func stringifySliceData(val reflect.Value, level uint8, ts *typeSpec) (string, error) {
	if len(ts.dimensions) <= int(level) {
		if val.Len() == 0 {
			return "", errors.New("stringify: incomplete value, failed to detect type")
		}
		ts.dimensions = append(ts.dimensions, val.Len())
	} else if ts.dimensions[level] != val.Len() {
		return "", errors.New("stringify: multidimensional slices must be rectangular")
	}

	if ctype, err := ctypes.FromType(val.Type().Elem()); err == nil {
		ts.ctype = ctype
		values := []string{}
		for idx := 0; idx < val.Len(); idx++ {
			if s, err := ctypes.ToString(ts.ctype, val.Index(idx).Interface(), true); err == nil {
				values = append(values, s)
			}
		}
		return dumpValues(values, level), nil
	}

	rv := strings.Builder{}
	rv.WriteString(lpadding(level) + "{")
	if val.Len() > 0 {
		rv.WriteString("\n")
	}

	if val.Len() > 0 {
		for idx := 0; idx < val.Len(); idx++ {
			d, err := stringify(val.Index(idx).Interface(), level+1, ts)
			if err != nil {
				return "", err
			}
			rv.WriteString(d + ",\n")
		}
	}

	if val.Len() > 0 {
		rv.WriteString(lpadding(level))
	}
	return rv.String() + "}", nil
}
