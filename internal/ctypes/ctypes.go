package ctypes

import (
	"fmt"
	"reflect"
)

type ctype struct {
	name      string
	cname     string
	isNumeric bool
	kind      reflect.Kind
	typ       reflect.Type
	toString  func(i interface{}, hex bool) string
}

var ctypes = []ctype{
	{"bool", "bool", false, reflect.Bool, reflect.TypeOf(false), printformat},
	{"int", "int32_t", true, reflect.Int, reflect.TypeOf(int(0)), valueformat},
	{"int8", "int8_t", true, reflect.Int8, reflect.TypeOf(int8(0)), valueformat},
	{"int16", "int16_t", true, reflect.Int16, reflect.TypeOf(int16(0)), valueformat},
	{"int32", "int32_t", true, reflect.Int32, reflect.TypeOf(int32(0)), valueformat},
	{"int64", "int64_t", true, reflect.Int64, reflect.TypeOf(int64(0)), valueformat},
	{"uint", "uint32_t", true, reflect.Uint, reflect.TypeOf(uint(0)), valueformat},
	{"uint8", "uint8_t", true, reflect.Uint8, reflect.TypeOf(uint8(0)), valueformat},
	{"uint16", "uint16_t", true, reflect.Uint16, reflect.TypeOf(uint16(0)), valueformat},
	{"uint32", "uint32_t", true, reflect.Uint32, reflect.TypeOf(uint32(0)), valueformat},
	{"uint64", "uint64_t", true, reflect.Uint64, reflect.TypeOf(uint64(0)), valueformat},
	{"float32", "float", true, reflect.Float32, reflect.TypeOf(float32(0)), printformat},
	{"float64", "double", true, reflect.Float64, reflect.TypeOf(float64(0)), printformat},
	{"string", "char*", false, reflect.String, reflect.TypeOf(""), stringformat},
}

func Validate(name string) bool {
	for _, ct := range ctypes {
		if ct.name == name || ct.cname == name {
			return true
		}
	}
	return false
}

func TypeIsNumeric(typ reflect.Type) bool {
	for _, ct := range ctypes {
		if ct.kind == typ.Kind() {
			return ct.isNumeric
		}
	}
	return false
}

func ToString(name string, obj interface{}, hex bool) (string, error) {
	for _, ct := range ctypes {
		if ct.name == name || ct.cname == name {
			return ct.toString(obj, hex), nil
		}
	}
	return "", fmt.Errorf("ctypes: type not supported: %s", name)
}

func ToType(name string) (reflect.Type, error) {
	for _, ct := range ctypes {
		if ct.name == name || ct.cname == name {
			return ct.typ, nil
		}
	}
	return nil, fmt.Errorf("ctypes: type not supported: %s", name)
}

func FromKind(kind reflect.Kind) (string, error) {
	for _, ct := range ctypes {
		if ct.kind == kind {
			return ct.cname, nil
		}
	}
	return "", fmt.Errorf("ctypes: kind not supported: %s", kind)
}
