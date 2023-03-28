package stringify

import (
	"testing"
)

type s1 struct {
	Foo string
}

type s2 struct {
	s1
	BarLol string
}

var stringifyArgs = []struct {
	itf          interface{}
	expectedData string
	expectedType string
	expectedDim  []int
}{
	{true, "true", "bool", []int{}},
	{int(123), "0x0000007b", "int32_t", []int{}},
	{int8(123), "0x7b", "int8_t", []int{}},
	{int16(123), "0x007b", "int16_t", []int{}},
	{int32(123), "0x0000007b", "int32_t", []int{}},
	{int64(123), "0x000000000000007b", "int64_t", []int{}},
	{int(-123), "0xffffff85", "int32_t", []int{}},
	{int8(-123), "0x85", "int8_t", []int{}},
	{int16(-123), "0xff85", "int16_t", []int{}},
	{int32(-123), "0xffffff85", "int32_t", []int{}},
	{int64(-123), "0xffffffffffffff85", "int64_t", []int{}},
	{uint(123), "0x0000007b", "uint32_t", []int{}},
	{uint8(123), "0x7b", "uint8_t", []int{}},
	{uint16(123), "0x007b", "uint16_t", []int{}},
	{uint32(123), "0x0000007b", "uint32_t", []int{}},
	{uint64(123), "0x000000000000007b", "uint64_t", []int{}},
	{float32(123.0), "123", "float", []int{}},
	{float32(123.4), "123.4", "float", []int{}},
	{float32(-123.0), "-123", "float", []int{}},
	{float32(-123.4), "-123.4", "float", []int{}},
	{float64(123.0), "123", "double", []int{}},
	{float64(123.4), "123.4", "double", []int{}},
	{float64(-123.0), "-123", "double", []int{}},
	{float64(-123.4), "-123.4", "double", []int{}},
	{"", "\"\"", "char*", []int{}},
	{"asd", "\"asd\"", "char*", []int{}},
	{struct{}{}, "{}", "struct {}", []int{}},
	{s1{}, "{\n    \"\",\n}", "struct {\n    char* foo;\n}", []int{}},
	{s1{"a"}, "{\n    \"a\",\n}", "struct {\n    char* foo;\n}", []int{}},
	{
		s2{s1: s1{"a"}, BarLol: "b"},
		"{\n    \"a\", \"b\",\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{},
	},
	{
		[]s1{{}},
		"{\n    {\n        \"\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1},
	},
	{
		[]s1{{"a"}},
		"{\n    {\n        \"a\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1},
	},
	{
		[]s1{{"a"}, {"b"}},
		"{\n    {\n        \"a\",\n    },\n    {\n        \"b\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{2},
	},
	{
		[]s2{{}},
		"{\n    {\n        \"\", \"\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1},
	},
	{
		[]s2{{s1: s1{"a"}, BarLol: "b"}},
		"{\n    {\n        \"a\", \"b\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1},
	},
	{
		[]s2{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}},
		"{\n    {\n        \"a\", \"b\",\n    },\n    {\n        \"c\", \"d\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{2},
	},
	{
		[][]s1{{{}}},
		"{\n    {\n        {\n            \"\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 1},
	},
	{
		[][]s1{{{"a"}}},
		"{\n    {\n        {\n            \"a\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 1},
	},
	{
		[][]s1{{{"a"}, {"b"}}},
		"{\n    {\n        {\n            \"a\",\n        },\n        {\n            \"b\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 2},
	},
	{
		[][]s2{{{}}},
		"{\n    {\n        {\n            \"\", \"\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1},
	},
	{
		[][]s2{{{s1: s1{"a"}, BarLol: "b"}}},
		"{\n    {\n        {\n            \"a\", \"b\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1},
	},
	{
		[][]s2{{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}}},
		"{\n    {\n        {\n            \"a\", \"b\",\n        },\n        {\n            \"c\", \"d\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 2},
	},
	{
		[][][]s2{{{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}}}},
		"{\n    {\n        {\n            {\n                \"a\", \"b\",\n            },\n            {\n                \"c\", \"d\",\n            },\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1, 2},
	},
	{
		[]int{1},
		"{\n    0x00000001,\n}",
		"int32_t",
		[]int{1},
	},
	{
		[]int{1, 2},
		"{\n    0x00000001, 0x00000002,\n}",
		"int32_t",
		[]int{2},
	},
	{
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		"{\n    0x00000001, 0x00000002, 0x00000003, 0x00000004, 0x00000005, 0x00000006, 0x00000007, 0x00000008,\n    0x00000009, 0x0000000a, 0x0000000b, 0x0000000c, 0x0000000d, 0x0000000e, 0x0000000f, 0x00000010,\n}",
		"int32_t",
		[]int{16},
	},
	{
		[][]int{{1, 2}},
		"{\n    {\n        0x00000001, 0x00000002,\n    },\n}",
		"int32_t",
		[]int{1, 2},
	},
	{
		[][]int{{1, 2}, {3, 4}},
		"{\n    {\n        0x00000001, 0x00000002,\n    },\n    {\n        0x00000003, 0x00000004,\n    },\n}",
		"int32_t",
		[]int{2, 2},
	},
	{
		[][]string{{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}},
		"{\n    {\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n    },\n}",
		"char*",
		[]int{1, 1},
	},
	{
		[][]string{{"c", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}},
		"{\n    {\n        \"c\",\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n    },\n}",
		"char*",
		[]int{1, 2},
	},
	{
		[][]string{{"c", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "d"}},
		"{\n    {\n        \"c\",\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n        \"d\",\n    },\n}",
		"char*",
		[]int{1, 3},
	},
	{
		[][]string{{"asd", "qwe"}, {"xcv", "vbn"}},
		"{\n    {\n        \"asd\", \"qwe\",\n    },\n    {\n        \"xcv\", \"vbn\",\n    },\n}",
		"char*",
		[]int{2, 2},
	},
}

func TestStringify(t *testing.T) {
	for _, tt := range stringifyArgs {
		data, ctype, dim, err := Stringify(tt.itf)
		if err != nil {
			t.Error(err)
		}

		if data != tt.expectedData {
			t.Errorf("expected data: got %q, want %q", data, tt.expectedData)
		}

		if ctype != tt.expectedType {
			t.Errorf("expected type: got %q, want %q", ctype, tt.expectedType)
		}

		if !func() bool {
			if len(dim) != len(tt.expectedDim) {
				return false
			}
			for i := range dim {
				if dim[i] != tt.expectedDim[i] {
					return false
				}
			}
			return true
		}() {
			t.Errorf("expected dimensions: got %v, want %v", dim, tt.expectedDim)
		}
	}
}

var stringifyValueArgs = []struct {
	itf          interface{}
	hex          bool
	expectedData string
}{
	{true, false, "true"},
	{int(123), false, "123"},
	{int8(123), false, "123"},
	{int16(123), false, "123"},
	{int32(123), false, "123"},
	{int64(123), false, "123"},
	{int(-123), false, "-123"},
	{int8(-123), false, "-123"},
	{int16(-123), false, "-123"},
	{int32(-123), false, "-123"},
	{int64(-123), false, "-123"},
	{uint(123), false, "123"},
	{uint8(123), false, "123"},
	{uint16(123), false, "123"},
	{uint32(123), false, "123"},
	{uint64(123), false, "123"},
	{float32(123.0), false, "123"},
	{float32(123.4), false, "123.4"},
	{float32(-123.0), false, "-123"},
	{float32(-123.4), false, "-123.4"},
	{float64(123.0), false, "123"},
	{float64(123.4), false, "123.4"},
	{float64(-123.0), false, "-123"},
	{float64(-123.4), false, "-123.4"},
	{"", false, "\"\""},
	{"asd", false, "\"asd\""},
	{true, true, "true"},
	{int(123), true, "0x0000007b"},
	{int8(123), true, "0x7b"},
	{int16(123), true, "0x007b"},
	{int32(123), true, "0x0000007b"},
	{int64(123), true, "0x000000000000007b"},
	{int(-123), true, "0xffffff85"},
	{int8(-123), true, "0x85"},
	{int16(-123), true, "0xff85"},
	{int32(-123), true, "0xffffff85"},
	{int64(-123), true, "0xffffffffffffff85"},
	{uint(123), true, "0x0000007b"},
	{uint8(123), true, "0x7b"},
	{uint16(123), true, "0x007b"},
	{uint32(123), true, "0x0000007b"},
	{uint64(123), true, "0x000000000000007b"},
	{float32(123.0), true, "123"},
	{float32(123.4), true, "123.4"},
	{float32(-123.0), true, "-123"},
	{float32(-123.4), true, "-123.4"},
	{float64(123.0), true, "123"},
	{float64(123.4), true, "123.4"},
	{float64(-123.0), true, "-123"},
	{float64(-123.4), true, "-123.4"},
	{"", true, "\"\""},
	{"asd", true, "\"asd\""},
}

func TestStringifyValue(t *testing.T) {
	for _, tt := range stringifyValueArgs {
		data, err := StringifyValue(tt.itf, tt.hex)
		if err != nil {
			t.Error(err)
		}

		if data != tt.expectedData {
			t.Errorf("expected data: got %q, want %q", data, tt.expectedData)
		}
	}
}

var stringifyErrorArgs = []struct {
	itf         interface{}
	expectedErr string
}{
	{nil, "stringify: got nil"},
	{complex(1, 1), "stringify: invalid type"},
	{[]int{}, "stringify: incomplete value, failed to detect type"},
	{[][]int{}, "stringify: incomplete value, failed to detect type"},
	{[][]int{{}}, "stringify: incomplete value, failed to detect type"},
	{[][]int{{1, 2}, {1}}, "stringify: multidimensional slices must be rectangular"},
}

func TestStringifyError(t *testing.T) {
	for _, tt := range stringifyErrorArgs {
		_, _, _, err := Stringify(tt.itf)
		if err == nil {
			t.Error("no error")
			continue
		}

		if err.Error() != tt.expectedErr {
			t.Errorf("expected error: got %q, want %q", err, tt.expectedErr)
		}
	}
}

var stringifyValueErrorArgs = []struct {
	itf         interface{}
	hex         bool
	expectedErr string
}{
	{nil, false, "stringify: got nil"},
	{complex(1, 1), false, "stringify: invalid type"},
	{nil, true, "stringify: got nil"},
	{complex(1, 1), true, "stringify: invalid type"},
}

func TestStringifyValueError(t *testing.T) {
	for _, tt := range stringifyValueErrorArgs {
		_, err := StringifyValue(tt.itf, tt.hex)
		if err == nil {
			t.Error("no error")
			continue
		}

		if err.Error() != tt.expectedErr {
			t.Errorf("expected error: got %q, want %q", err, tt.expectedErr)
		}
	}
}
