package convert

import (
	"errors"
	"fmt"
	"reflect"

	"rafaelmartins.com/p/synth-datagen/internal/ctypes"
)

func Slice(slice interface{}, to string) (interface{}, error) {
	if slice == nil {
		return nil, errors.New("slice: got nil")
	}

	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return nil, errors.New("slice: not a slice")
	}

	etype := val.Type().Elem()
	if etype.Kind() == reflect.Interface {
		if val.Len() == 0 {
			return nil, fmt.Errorf("slice: empty, can't guess type")
		}
		etype = reflect.TypeOf(val.Index(0).Interface())
	}
	if etype == nil {
		return nil, fmt.Errorf("slice: unsupported element type: nil")
	}
	if etype.Kind() != reflect.Slice && !ctypes.TypeIsScalar(etype) {
		return nil, fmt.Errorf("slice: unsupported element type: %s", etype)
	}

	typ := etype
	if to != "" {
		var err error
		typ, err = ctypes.ToType(to)
		if err != nil {
			return nil, err
		}
	}

	rv := reflect.Value{}
	for i := 0; i < val.Len(); i++ {
		eval := val.Index(i)
		if eval.Kind() == reflect.Interface {
			eval = reflect.ValueOf(eval.Interface())
		}
		if eval.Kind() == reflect.Slice {
			r, err := Slice(eval.Interface(), to)
			if err != nil {
				return nil, err
			}
			if rv.Kind() == reflect.Invalid {
				rv = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(r)), 0, val.Len())
			}
			rv = reflect.Append(rv, reflect.ValueOf(r))
			continue
		}
		if !eval.CanConvert(typ) {
			return nil, fmt.Errorf("slice: value of type %s cannot be converted to type %s", eval.Type(), typ)
		}
		if rv.Kind() == reflect.Invalid {
			rv = reflect.MakeSlice(reflect.SliceOf(typ), 0, val.Len())
		}
		rv = reflect.Append(rv, eval.Convert(typ))
	}
	if rv.Kind() == reflect.Invalid {
		return nil, nil
	}
	return rv.Interface(), nil
}
