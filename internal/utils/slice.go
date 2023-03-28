package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rafaelmartins/synth-datagen/internal/ctypes"
)

func CastFloatSliceTo(s []float64, t string) (interface{}, error) {
	typ, err := ctypes.ToType(t)
	if err != nil {
		return nil, err
	}
	rv := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(s))
	for _, v := range s {
		val := reflect.ValueOf(v)
		if !val.CanConvert(typ) {
			return nil, fmt.Errorf("slice: value of type %s cannot be converted to type %s", val.Type(), typ)
		}
		rv = reflect.Append(rv, val.Convert(typ))
	}
	return rv.Interface(), nil
}

func CastInterfaceSlice(s []interface{}) (interface{}, error) {
	if len(s) == 0 {
		return nil, errors.New("slice: empty, can't guess type")
	}
	typ := reflect.TypeOf(s[0])
	rv := reflect.MakeSlice(reflect.SliceOf(typ), 0, len(s))
	for _, v := range s {
		val := reflect.ValueOf(v)
		if !val.CanConvert(typ) {
			return nil, fmt.Errorf("slice: value of type %s cannot be converted to type %s", val.Type(), typ)
		}
		rv = reflect.Append(rv, val.Convert(typ))
	}
	return rv.Interface(), nil
}
