package stringify

import (
	"errors"
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
)

type typeSpec struct {
	ctype      string
	stype      []*structSpec
	dimensions []int
}

func stringify(obj interface{}, level uint8, ts *typeSpec) (string, error) {
	if obj == nil {
		return "", errors.New("stringify: got nil")
	}

	val := reflect.ValueOf(obj)
	if ctype, err := ctypes.FromType(val.Type()); err == nil {
		ts.ctype = ctype
		return ctypes.ToString(ts.ctype, obj, true)
	}

	switch val.Kind() {
	case reflect.Struct:
		return stringifyStructData(val, level, ts), nil
	case reflect.Slice:
		return stringifySliceData(val, level, ts)
	default:
		return "", errors.New("stringify: invalid type")
	}
}

func Stringify(obj interface{}) (string, string, []int, error) {
	ts := &typeSpec{}
	data, err := stringify(obj, 0, ts)
	if err != nil {
		return "", "", []int{}, err
	}

	ctype := ts.ctype
	if ctype == "" {
		ctype = stringifyStructType(ts)
	}

	return data, ctype, ts.dimensions, nil
}

func StringifyValue(obj interface{}, hex bool) (string, error) {
	if obj == nil {
		return "", errors.New("stringify: got nil")
	}
	if ctype, err := ctypes.FromType(reflect.TypeOf(obj)); err == nil {
		return ctypes.ToString(ctype, obj, hex)
	}
	return "", errors.New("stringify: invalid type")
}
