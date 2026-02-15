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

type s3 struct {
	Name  string
	Count int32
	Flag  bool
	Rate  float64
}

type s4 struct {
	Name   string
	hidden int32
}

type s5 struct {
	Name string
	C    complex128
	Flag bool
}

var stringifyArgs = []struct {
	name         string
	itf          interface{}
	expectedData string
	expectedType string
	expectedDim  []int
}{
	{"bool_true", true, "true", "bool", []int{}},
	{"int_pos", int(123), "0x0000007b", "int32_t", []int{}},
	{"int8_pos", int8(123), "0x7b", "int8_t", []int{}},
	{"int16_pos", int16(123), "0x007b", "int16_t", []int{}},
	{"int32_pos", int32(123), "0x0000007b", "int32_t", []int{}},
	{"int64_pos", int64(123), "0x000000000000007b", "int64_t", []int{}},
	{"int_neg", int(-123), "0xffffff85", "int32_t", []int{}},
	{"int8_neg", int8(-123), "0x85", "int8_t", []int{}},
	{"int16_neg", int16(-123), "0xff85", "int16_t", []int{}},
	{"int32_neg", int32(-123), "0xffffff85", "int32_t", []int{}},
	{"int64_neg", int64(-123), "0xffffffffffffff85", "int64_t", []int{}},
	{"uint_pos", uint(123), "0x0000007b", "uint32_t", []int{}},
	{"uint8_pos", uint8(123), "0x7b", "uint8_t", []int{}},
	{"uint16_pos", uint16(123), "0x007b", "uint16_t", []int{}},
	{"uint32_pos", uint32(123), "0x0000007b", "uint32_t", []int{}},
	{"uint64_pos", uint64(123), "0x000000000000007b", "uint64_t", []int{}},
	{"float32_int", float32(123.0), "123", "float", []int{}},
	{"float32_frac", float32(123.4), "123.4", "float", []int{}},
	{"float32_neg_int", float32(-123.0), "-123", "float", []int{}},
	{"float32_neg_frac", float32(-123.4), "-123.4", "float", []int{}},
	{"float64_int", float64(123.0), "123", "double", []int{}},
	{"float64_frac", float64(123.4), "123.4", "double", []int{}},
	{"float64_neg_int", float64(-123.0), "-123", "double", []int{}},
	{"float64_neg_frac", float64(-123.4), "-123.4", "double", []int{}},
	{"string_empty", "", "\"\"", "char*", []int{}},
	{"string_value", "asd", "\"asd\"", "char*", []int{}},
	{"empty_struct", struct{}{}, "{}", "struct {}", []int{}},
	{"s1_zero", s1{}, "{\n    \"\",\n}", "struct {\n    char* foo;\n}", []int{}},
	{"s1_value", s1{"a"}, "{\n    \"a\",\n}", "struct {\n    char* foo;\n}", []int{}},
	{
		"s2_embedded",
		s2{s1: s1{"a"}, BarLol: "b"},
		"{\n    \"a\", \"b\",\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{},
	},
	{
		"slice_s1_one_zero",
		[]s1{{}},
		"{\n    {\n        \"\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1},
	},
	{
		"slice_s1_one",
		[]s1{{"a"}},
		"{\n    {\n        \"a\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1},
	},
	{
		"slice_s1_two",
		[]s1{{"a"}, {"b"}},
		"{\n    {\n        \"a\",\n    },\n    {\n        \"b\",\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{2},
	},
	{
		"slice_s2_one_zero",
		[]s2{{}},
		"{\n    {\n        \"\", \"\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1},
	},
	{
		"slice_s2_one",
		[]s2{{s1: s1{"a"}, BarLol: "b"}},
		"{\n    {\n        \"a\", \"b\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1},
	},
	{
		"slice_s2_two",
		[]s2{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}},
		"{\n    {\n        \"a\", \"b\",\n    },\n    {\n        \"c\", \"d\",\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{2},
	},
	{
		"2d_s1_1x1_zero",
		[][]s1{{{}}},
		"{\n    {\n        {\n            \"\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 1},
	},
	{
		"2d_s1_1x1",
		[][]s1{{{"a"}}},
		"{\n    {\n        {\n            \"a\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 1},
	},
	{
		"2d_s1_1x2",
		[][]s1{{{"a"}, {"b"}}},
		"{\n    {\n        {\n            \"a\",\n        },\n        {\n            \"b\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n}",
		[]int{1, 2},
	},
	{
		"2d_s2_1x1_zero",
		[][]s2{{{}}},
		"{\n    {\n        {\n            \"\", \"\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1},
	},
	{
		"2d_s2_1x1",
		[][]s2{{{s1: s1{"a"}, BarLol: "b"}}},
		"{\n    {\n        {\n            \"a\", \"b\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1},
	},
	{
		"2d_s2_1x2",
		[][]s2{{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}}},
		"{\n    {\n        {\n            \"a\", \"b\",\n        },\n        {\n            \"c\", \"d\",\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 2},
	},
	{
		"3d_s2_1x1x2",
		[][][]s2{{{{s1: s1{"a"}, BarLol: "b"}, {s1: s1{"c"}, BarLol: "d"}}}},
		"{\n    {\n        {\n            {\n                \"a\", \"b\",\n            },\n            {\n                \"c\", \"d\",\n            },\n        },\n    },\n}",
		"struct {\n    char* foo;\n    char* bar_lol;\n}",
		[]int{1, 1, 2},
	},
	{
		"slice_int_one",
		[]int{1},
		"{\n    0x00000001,\n}",
		"int32_t",
		[]int{1},
	},
	{
		"slice_int_two",
		[]int{1, 2},
		"{\n    0x00000001, 0x00000002,\n}",
		"int32_t",
		[]int{2},
	},
	{
		"slice_int_wrapping",
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		"{\n    0x00000001, 0x00000002, 0x00000003, 0x00000004, 0x00000005, 0x00000006, 0x00000007, 0x00000008,\n    0x00000009, 0x0000000a, 0x0000000b, 0x0000000c, 0x0000000d, 0x0000000e, 0x0000000f, 0x00000010,\n}",
		"int32_t",
		[]int{16},
	},
	{
		"2d_int_1x2",
		[][]int{{1, 2}},
		"{\n    {\n        0x00000001, 0x00000002,\n    },\n}",
		"int32_t",
		[]int{1, 2},
	},
	{
		"2d_int_2x2",
		[][]int{{1, 2}, {3, 4}},
		"{\n    {\n        0x00000001, 0x00000002,\n    },\n    {\n        0x00000003, 0x00000004,\n    },\n}",
		"int32_t",
		[]int{2, 2},
	},
	{
		"2d_string_long_single",
		[][]string{{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}},
		"{\n    {\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n    },\n}",
		"char*",
		[]int{1, 1},
	},
	{
		"2d_string_long_with_short",
		[][]string{{"c", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}},
		"{\n    {\n        \"c\",\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n    },\n}",
		"char*",
		[]int{1, 2},
	},
	{
		"2d_string_long_three",
		[][]string{{"c", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "d"}},
		"{\n    {\n        \"c\",\n        \"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\",\n        \"d\",\n    },\n}",
		"char*",
		[]int{1, 3},
	},
	{
		"2d_string_2x2",
		[][]string{{"asd", "qwe"}, {"xcv", "vbn"}},
		"{\n    {\n        \"asd\", \"qwe\",\n    },\n    {\n        \"xcv\", \"vbn\",\n    },\n}",
		"char*",
		[]int{2, 2},
	},
	{"bool_false", false, "false", "bool", []int{}},
	{"int_zero", int(0), "0x00000000", "int32_t", []int{}},
	{"uint_zero", uint(0), "0x00000000", "uint32_t", []int{}},
	{"float32_zero", float32(0), "0", "float", []int{}},
	{"float64_zero", float64(0), "0", "double", []int{}},
	{
		"struct_mixed_types",
		s3{Name: "x", Count: 5, Flag: true, Rate: 1.5},
		"{\n    \"x\", 0x00000005, true, 1.5,\n}",
		"struct {\n    char* name;\n    int32_t count;\n    bool flag;\n    double rate;\n}",
		[]int{},
	},
	{
		"struct_unexported_field",
		s4{Name: "x", hidden: 42},
		"{\n    \"x\",\n}",
		"struct {\n    char* name;\n}",
		[]int{},
	},
	{
		"struct_unsupported_field",
		s5{Name: "x", Flag: true},
		"{\n    \"x\", true,\n}",
		"struct {\n    char* name;\n    bool flag;\n}",
		[]int{},
	},
	{"slice_uint8", []uint8{0x0a, 0x0b}, "{\n    0x0a, 0x0b,\n}", "uint8_t", []int{2}},
	{"slice_float32", []float32{1.5, 2.5}, "{\n    1.5, 2.5,\n}", "float", []int{2}},
	{"slice_bool", []bool{true, false}, "{\n    true, false,\n}", "bool", []int{2}},
	{"slice_string", []string{"foo", "bar"}, "{\n    \"foo\", \"bar\",\n}", "char*", []int{2}},
	{
		"slice_s3_mixed",
		[]s3{{Name: "a", Count: 1, Flag: true, Rate: 0.5}},
		"{\n    {\n        \"a\", 0x00000001, true, 0.5,\n    },\n}",
		"struct {\n    char* name;\n    int32_t count;\n    bool flag;\n    double rate;\n}",
		[]int{1},
	},
}

func TestStringify(t *testing.T) {
	for _, tt := range stringifyArgs {
		t.Run(tt.name, func(t *testing.T) {
			data, ctype, dim, err := Stringify(tt.itf)
			if err != nil {
				t.Fatal(err)
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
		})
	}
}

var stringifyValueArgs = []struct {
	name         string
	itf          interface{}
	hex          bool
	expectedData string
}{
	{"dec_bool_true", true, false, "true"},
	{"dec_int_pos", int(123), false, "123"},
	{"dec_int8_pos", int8(123), false, "123"},
	{"dec_int16_pos", int16(123), false, "123"},
	{"dec_int32_pos", int32(123), false, "123"},
	{"dec_int64_pos", int64(123), false, "123"},
	{"dec_int_neg", int(-123), false, "-123"},
	{"dec_int8_neg", int8(-123), false, "-123"},
	{"dec_int16_neg", int16(-123), false, "-123"},
	{"dec_int32_neg", int32(-123), false, "-123"},
	{"dec_int64_neg", int64(-123), false, "-123"},
	{"dec_uint_pos", uint(123), false, "123"},
	{"dec_uint8_pos", uint8(123), false, "123"},
	{"dec_uint16_pos", uint16(123), false, "123"},
	{"dec_uint32_pos", uint32(123), false, "123"},
	{"dec_uint64_pos", uint64(123), false, "123"},
	{"dec_float32_int", float32(123.0), false, "123"},
	{"dec_float32_frac", float32(123.4), false, "123.4"},
	{"dec_float32_neg_int", float32(-123.0), false, "-123"},
	{"dec_float32_neg_frac", float32(-123.4), false, "-123.4"},
	{"dec_float64_int", float64(123.0), false, "123"},
	{"dec_float64_frac", float64(123.4), false, "123.4"},
	{"dec_float64_neg_int", float64(-123.0), false, "-123"},
	{"dec_float64_neg_frac", float64(-123.4), false, "-123.4"},
	{"dec_string_empty", "", false, "\"\""},
	{"dec_string_value", "asd", false, "\"asd\""},
	{"hex_bool_true", true, true, "true"},
	{"hex_int_pos", int(123), true, "0x0000007b"},
	{"hex_int8_pos", int8(123), true, "0x7b"},
	{"hex_int16_pos", int16(123), true, "0x007b"},
	{"hex_int32_pos", int32(123), true, "0x0000007b"},
	{"hex_int64_pos", int64(123), true, "0x000000000000007b"},
	{"hex_int_neg", int(-123), true, "0xffffff85"},
	{"hex_int8_neg", int8(-123), true, "0x85"},
	{"hex_int16_neg", int16(-123), true, "0xff85"},
	{"hex_int32_neg", int32(-123), true, "0xffffff85"},
	{"hex_int64_neg", int64(-123), true, "0xffffffffffffff85"},
	{"hex_uint_pos", uint(123), true, "0x0000007b"},
	{"hex_uint8_pos", uint8(123), true, "0x7b"},
	{"hex_uint16_pos", uint16(123), true, "0x007b"},
	{"hex_uint32_pos", uint32(123), true, "0x0000007b"},
	{"hex_uint64_pos", uint64(123), true, "0x000000000000007b"},
	{"hex_float32_int", float32(123.0), true, "123"},
	{"hex_float32_frac", float32(123.4), true, "123.4"},
	{"hex_float32_neg_int", float32(-123.0), true, "-123"},
	{"hex_float32_neg_frac", float32(-123.4), true, "-123.4"},
	{"hex_float64_int", float64(123.0), true, "123"},
	{"hex_float64_frac", float64(123.4), true, "123.4"},
	{"hex_float64_neg_int", float64(-123.0), true, "-123"},
	{"hex_float64_neg_frac", float64(-123.4), true, "-123.4"},
	{"hex_string_empty", "", true, "\"\""},
	{"hex_string_value", "asd", true, "\"asd\""},
	{"dec_bool_false", false, false, "false"},
	{"dec_int_zero", int(0), false, "0"},
	{"dec_uint_zero", uint(0), false, "0"},
	{"dec_float32_zero", float32(0), false, "0"},
	{"dec_float64_zero", float64(0), false, "0"},
	{"hex_bool_false", false, true, "false"},
	{"hex_int_zero", int(0), true, "0x00000000"},
	{"hex_uint_zero", uint(0), true, "0x00000000"},
	{"hex_float32_zero", float32(0), true, "0"},
	{"hex_float64_zero", float64(0), true, "0"},
}

func TestStringifyValue(t *testing.T) {
	for _, tt := range stringifyValueArgs {
		t.Run(tt.name, func(t *testing.T) {
			data, err := StringifyValue(tt.itf, tt.hex)
			if err != nil {
				t.Fatal(err)
			}

			if data != tt.expectedData {
				t.Errorf("expected data: got %q, want %q", data, tt.expectedData)
			}
		})
	}
}

var stringifyErrorArgs = []struct {
	name        string
	itf         interface{}
	expectedErr string
}{
	{"nil", nil, "stringify: got nil"},
	{"complex", complex(1, 1), "stringify: invalid type"},
	{"empty_slice", []int{}, "stringify: incomplete value, failed to detect type"},
	{"empty_2d_slice", [][]int{}, "stringify: incomplete value, failed to detect type"},
	{"empty_inner_2d_slice", [][]int{{}}, "stringify: incomplete value, failed to detect type"},
	{"non_rectangular", [][]int{{1, 2}, {1}}, "stringify: multidimensional slices must be rectangular"},
	{"slice_complex", []complex128{1}, "stringify: invalid type"},
}

func TestStringifyError(t *testing.T) {
	for _, tt := range stringifyErrorArgs {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := Stringify(tt.itf)
			if err == nil {
				t.Fatal("no error")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error: got %q, want %q", err, tt.expectedErr)
			}
		})
	}
}

var stringifyValueErrorArgs = []struct {
	name        string
	itf         interface{}
	hex         bool
	expectedErr string
}{
	{"dec_nil", nil, false, "stringify: got nil"},
	{"dec_complex", complex(1, 1), false, "stringify: invalid type"},
	{"hex_nil", nil, true, "stringify: got nil"},
	{"hex_complex", complex(1, 1), true, "stringify: invalid type"},
	{"dec_struct", struct{}{}, false, "stringify: invalid type"},
	{"dec_slice", []int{1}, false, "stringify: invalid type"},
	{"hex_struct", struct{}{}, true, "stringify: invalid type"},
	{"hex_slice", []int{1}, true, "stringify: invalid type"},
}

func TestStringifyValueError(t *testing.T) {
	for _, tt := range stringifyValueErrorArgs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StringifyValue(tt.itf, tt.hex)
			if err == nil {
				t.Fatal("no error")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error: got %q, want %q", err, tt.expectedErr)
			}
		})
	}
}
