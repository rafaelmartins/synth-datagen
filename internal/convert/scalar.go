package convert

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
)

func Scalar(scalar interface{}, to string) (interface{}, error) {
	if scalar == nil {
		return nil, errors.New("scalar: got nil")
	}
	if to == "" {
		return nil, errors.New("scalar: missing type")
	}

	val := reflect.ValueOf(scalar)
	if !ctypes.TypeIsScalar(val.Type()) {
		return nil, errors.New("scalar: not a scalar")
	}

	typ, err := ctypes.ToType(to)
	if err != nil {
		return nil, err
	}
	if !val.CanConvert(typ) {
		return nil, fmt.Errorf("slice: value of type %s cannot be converted to type %s", val.Type(), typ)
	}
	return val.Convert(typ).Interface(), nil
}
