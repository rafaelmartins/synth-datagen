package stringify

import (
	"errors"
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
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

	if ctype, err := ctypes.FromKind(val.Type().Elem().Kind()); err == nil {
		ts.ctype = ctype
		values := []string{}
		for idx := 0; idx < val.Len(); idx++ {
			if s, err := ctypes.ToString(ts.ctype, val.Index(idx).Interface(), true); err == nil {
				values = append(values, s)
			}
		}
		return dumpValues(values, level), nil
	}

	rv := lpadding(level) + "{"
	if val.Len() > 0 {
		rv += "\n"
	}

	if val.Len() > 0 {
		for idx := 0; idx < val.Len(); idx++ {
			d, err := stringify(val.Index(idx).Interface(), level+1, ts)
			if err != nil {
				return "", err
			}
			rv += d + ",\n"
		}
	}

	if val.Len() > 0 {
		rv += lpadding(level)
	}
	return rv + "}", nil
}
