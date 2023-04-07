package utils

import (
	"fmt"
	"io"
	"math"
	"os"
	"path"
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

func FormatStringSliceWidth(s []string, width int) ([]string, error) {
	rv := []string{}
	for _, v := range s {
		f := fmt.Sprintf("%*s", width, v)
		if m := int(math.Abs(float64(width))); len(f) != m {
			return nil, fmt.Errorf("width overflow: %q (%d > %d)", f, len(f), m)
		}
		rv = append(rv, f)
	}
	return rv, nil
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
