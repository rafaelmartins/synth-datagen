package codegen

import (
	"bytes"
	"strings"
	"testing"
)

func TestDataWrite(t *testing.T) {
	t.Run("scalar_int", func(t *testing.T) {
		var dl dataList
		dl.add("my_var", int32(42), nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\nstatic const int32_t my_var = 0x0000002a;\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("scalar_bool", func(t *testing.T) {
		var dl dataList
		dl.add("flag", true, nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\nstatic const bool flag = true;\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("scalar_string", func(t *testing.T) {
		var dl dataList
		dl.add("name", "hello", nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\nstatic const char* name = \"hello\";\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("1d_array_dim_defines", func(t *testing.T) {
		var dl dataList
		dl.add("arr", []int32{1, 2, 3}, nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const int32_t arr[3]") {
			t.Errorf("missing array declaration in %q", got)
		}
		if !strings.Contains(got, "#define arr_len 3\n") {
			t.Errorf("missing #define arr_len in %q", got)
		}
	})

	t.Run("2d_array_dim_defines", func(t *testing.T) {
		var dl dataList
		dl.add("mat", [][]int32{{1, 2}, {3, 4}}, nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const int32_t mat[2][2]") {
			t.Errorf("missing array declaration in %q", got)
		}
		if !strings.Contains(got, "#define mat_rows 2\n") {
			t.Errorf("missing #define mat_rows in %q", got)
		}
		if !strings.Contains(got, "#define mat_cols 2\n") {
			t.Errorf("missing #define mat_cols in %q", got)
		}
	})

	t.Run("3d_array_dim_defines", func(t *testing.T) {
		var dl dataList
		dl.add("cube", [][][]int32{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}, nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const int32_t cube[2][2][2]") {
			t.Errorf("missing array declaration in %q", got)
		}
		if !strings.Contains(got, "#define cube_len_0 2\n") {
			t.Errorf("missing #define cube_len_0 in %q", got)
		}
		if !strings.Contains(got, "#define cube_len_1 2\n") {
			t.Errorf("missing #define cube_len_1 in %q", got)
		}
		if !strings.Contains(got, "#define cube_len_2 2\n") {
			t.Errorf("missing #define cube_len_2 in %q", got)
		}
	})

	t.Run("attributes", func(t *testing.T) {
		var dl dataList
		dl.add("my_var", int32(1), []string{"__attribute__((aligned(4)))"}, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\nstatic const int32_t my_var __attribute__((aligned(4))) = 0x00000001;\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("multiple_attributes", func(t *testing.T) {
		var dl dataList
		dl.add("my_var", int32(1), []string{"__attr1__", "__attr2__"}, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\nstatic const int32_t my_var __attr1__ __attr2__ = 0x00000001;\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("attributes_with_array", func(t *testing.T) {
		var dl dataList
		dl.add("arr", []int32{1, 2}, []string{"__aligned__"}, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const int32_t arr[2] __aligned__ = {") {
			t.Errorf("unexpected output: %q", got)
		}
	})

	t.Run("str_width_string", func(t *testing.T) {
		var dl dataList
		dl.add("name", "hi", nil, new(5))
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const char name[5]") {
			t.Errorf("expected char array declaration, got %q", got)
		}
		if !strings.Contains(got, `"   hi"`) {
			t.Errorf("expected right-aligned string, got %q", got)
		}
		if !strings.Contains(got, "#define name_len 5\n") {
			t.Errorf("missing #define name_len in %q", got)
		}
	})

	t.Run("str_width_negative", func(t *testing.T) {
		var dl dataList
		dl.add("name", "hi", nil, new(-5))
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const char name[5]") {
			t.Errorf("expected char array with abs(width), got %q", got)
		}
		if !strings.Contains(got, `"hi   "`) {
			t.Errorf("expected left-aligned string, got %q", got)
		}
	})

	t.Run("str_width_string_slice", func(t *testing.T) {
		var dl dataList
		dl.add("names", []string{"hi", "bye"}, nil, new(5))
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "static const char names[2][5]") {
			t.Errorf("expected 2D char array, got %q", got)
		}
		if !strings.Contains(got, "#define names_rows 2\n") {
			t.Errorf("missing #define names_rows in %q", got)
		}
		if !strings.Contains(got, "#define names_cols 5\n") {
			t.Errorf("missing #define names_cols in %q", got)
		}
	})

	t.Run("str_width_overflow", func(t *testing.T) {
		var dl dataList
		dl.add("name", "toolong", nil, new(3))
		var buf bytes.Buffer
		err := dl.write(&buf)
		if err == nil {
			t.Fatal("expected width overflow error")
		}
		if !strings.Contains(err.Error(), "width overflow") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("multiple_data", func(t *testing.T) {
		var dl dataList
		dl.add("a", int32(1), nil, nil)
		dl.add("b", int32(2), nil, nil)
		var buf bytes.Buffer
		if err := dl.write(&buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		// each data entry starts with a newline
		if strings.Count(got, "static const") != 2 {
			t.Errorf("expected 2 data declarations, got %q", got)
		}
	})

	t.Run("error_nil_value", func(t *testing.T) {
		var dl dataList
		dl.add("bad", nil, nil, nil)
		var buf bytes.Buffer
		err := dl.write(&buf)
		if err == nil {
			t.Fatal("expected error for nil value")
		}
	})

	t.Run("error_invalid_type", func(t *testing.T) {
		var dl dataList
		dl.add("bad", complex(1, 2), nil, nil)
		var buf bytes.Buffer
		err := dl.write(&buf)
		if err == nil {
			t.Fatal("expected error for unsupported type")
		}
	})
}
