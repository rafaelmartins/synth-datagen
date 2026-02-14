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

func SliceStruct(slice interface{}, to string) (interface{}, error) {
	if slice == nil {
		return nil, errors.New("slicestruct: got nil")
	}

	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return nil, errors.New("slicestruct: not a slice")
	}

	etype := val.Type().Elem()
	if etype.Kind() == reflect.Interface {
		if val.Len() == 0 {
			return nil, fmt.Errorf("slicestruct: empty, can't guess type")
		}
		etype = reflect.TypeOf(val.Index(0).Interface())
	}
	if etype == nil {
		return nil, fmt.Errorf("slicestruct: unsupported element type: nil")
	}
	if etype.Kind() != reflect.Struct {
		return nil, fmt.Errorf("slicestruct: not a slice of structs")
	}

	fields := reflect.VisibleFields(etype)

	var typ reflect.Type
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		if typ == nil {
			typ = field.Type
			continue
		}

		if field.Type.Kind() != typ.Kind() {
			return nil, fmt.Errorf("slicestruct: all struct fields must have same type: %s != %s", field.Type, typ)
		}
	}

	if to != "" {
		var err error
		typ, err = ctypes.ToType(to)
		if err != nil {
			return nil, err
		}
	}

	nfields := []reflect.StructField{}
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		nfield := field
		nfield.Type = typ
		nfields = append(nfields, nfield)
	}
	rvt := reflect.StructOf(nfields)

	rv := reflect.Value{}
	for i := 0; i < val.Len(); i++ {
		eval := val.Index(i)
		if eval.Kind() == reflect.Interface {
			eval = reflect.ValueOf(eval.Interface())
		}
		if eval.Kind() != reflect.Struct {
			return nil, fmt.Errorf("slicestruct: unexpected type: %s", eval.Type())
		}
		if rv.Kind() == reflect.Invalid {
			rv = reflect.MakeSlice(reflect.SliceOf(rvt), 0, val.Len())
		}

		rvv := reflect.New(rvt).Elem()
		for j, nf := range nfields {
			eeval := eval.FieldByName(nf.Name)
			if !eeval.CanConvert(typ) {
				return nil, fmt.Errorf("slicestruct: value of type %s cannot be converted to type %s", eeval.Type(), typ)
			}
			rvv.Field(j).Set(eeval.Convert(typ))
		}
		rv = reflect.Append(rv, rvv)
	}

	if rv.Kind() == reflect.Invalid {
		return nil, nil
	}
	return rv.Interface(), nil
}
