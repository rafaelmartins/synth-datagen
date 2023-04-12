package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"unicode"
)

func FieldNameToSnake(name string) string {
	rv := ""
	for idx, c := range name {
		if idx == 0 {
			rv += string(unicode.ToLower(c))
			continue
		}

		if unicode.IsUpper(c) {
			rv += "_" + string(unicode.ToLower(c))
		} else if unicode.IsLower(c) || unicode.IsNumber(c) {
			rv += string(c)
		} else {
			rv += "_"
		}
	}
	return rv
}

func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func ApplyStringWidth(obj interface{}, width int) error {
	if obj == nil {
		return errors.New("got nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)

			if v.Kind() == reflect.Interface {
				v = reflect.ValueOf(v.Interface())
			}

			if v.Kind() == reflect.String && v.CanAddr() && v.Addr().CanInterface() {
				if err := ApplyStringWidth(v.Addr().Interface(), width); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if val.Kind() == reflect.Pointer {
		v := val.Elem()
		if v.Kind() == reflect.Interface {
			v = reflect.ValueOf(v.Interface())
		}

		if v.Kind() == reflect.String {
			s := fmt.Sprintf("%*s", width, v.String())
			if m := Abs(width); len(s) > m {
				return fmt.Errorf("width overflow: %q (%d > %d)", s, len(s), m)
			}
			val.Elem().Set(reflect.ValueOf(s))
		}
	}
	return nil
}

func WriteFile(name string, w interface{ Write(w io.Writer) error }) error {
	dir := path.Dir(name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	if w != nil {
		return w.Write(f)
	}
	return nil
}
