package ctypes

import (
	"fmt"
	"reflect"
)

type ctype struct {
	isNumeric bool
	zero      interface{}
	toString  func(i interface{}, hex bool) string
}

var ctypes = map[string]*ctype{
	"bool":     {false, false, printformat},
	"int8_t":   {true, int8(0), valueformat},
	"int16_t":  {true, int16(0), valueformat},
	"int32_t":  {true, int32(0), valueformat},
	"int64_t":  {true, int64(0), valueformat},
	"uint8_t":  {true, uint8(0), valueformat},
	"uint16_t": {true, uint16(0), valueformat},
	"uint32_t": {true, uint32(0), valueformat},
	"uint64_t": {true, uint64(0), valueformat},
	"float":    {true, float32(0), printformat},
	"double":   {true, float64(0), printformat},
	"char*":    {false, "", stringformat},
}

var kind2name = map[reflect.Kind]string{
	reflect.Bool:    "bool",
	reflect.Int:     "int32_t",
	reflect.Int8:    "int8_t",
	reflect.Int16:   "int16_t",
	reflect.Int32:   "int32_t",
	reflect.Int64:   "int64_t",
	reflect.Uint:    "uint32_t",
	reflect.Uint8:   "uint8_t",
	reflect.Uint16:  "uint16_t",
	reflect.Uint32:  "uint32_t",
	reflect.Uint64:  "uint64_t",
	reflect.Float32: "float",
	reflect.Float64: "double",
	reflect.String:  "char*",
}

func TypeIsScalar(typ reflect.Type) bool {
	_, found := kind2name[typ.Kind()]
	return found
}

func TypeIsNumeric(typ reflect.Type) bool {
	name, found := kind2name[typ.Kind()]
	if !found {
		return false
	}
	ct, found := ctypes[name]
	return found && ct.isNumeric
}

func ToString(name string, obj interface{}, hex bool) (string, error) {
	ct, found := ctypes[name]
	if !found {
		return "", fmt.Errorf("ctypes: type not supported: %s", name)
	}
	return ct.toString(obj, hex), nil
}

func ToType(name string) (reflect.Type, error) {
	ct, found := ctypes[name]
	if !found {
		return nil, fmt.Errorf("ctypes: type not supported: %s", name)
	}
	return reflect.TypeOf(ct.zero), nil
}

func FromType(typ reflect.Type) (string, error) {
	name, found := kind2name[typ.Kind()]
	if !found {
		return "", fmt.Errorf("ctypes: type not supported: %s", typ)
	}
	return name, nil
}
