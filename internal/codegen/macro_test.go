package codegen

import (
	"bytes"
	"testing"
)

func TestMacroWrite(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var ml macroList
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "" {
			t.Errorf("expected empty output, got %q", buf.String())
		}
	})

	t.Run("raw_string", func(t *testing.T) {
		var ml macroList
		ml.add("FOO", "bar", false, true)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define FOO bar\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("raw_int", func(t *testing.T) {
		var ml macroList
		ml.add("SIZE", 42, false, true)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define SIZE 42\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("raw_ignores_hex_flag", func(t *testing.T) {
		var ml macroList
		ml.add("SIZE", 42, true, true)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		// raw mode uses %v formatting, ignoring hex flag
		expected := "\n#define SIZE 42\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("stringify_decimal", func(t *testing.T) {
		var ml macroList
		ml.add("SIZE", int32(42), false, false)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define SIZE 42\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("stringify_hex", func(t *testing.T) {
		var ml macroList
		ml.add("SIZE", int32(42), true, false)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define SIZE 0x0000002a\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("stringify_string", func(t *testing.T) {
		var ml macroList
		ml.add("NAME", "hello", false, false)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define NAME \"hello\"\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("stringify_bool", func(t *testing.T) {
		var ml macroList
		ml.add("ENABLED", true, false, false)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define ENABLED true\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		var ml macroList
		ml.add("FOO", "bar", false, true)
		ml.add("SIZE", int32(42), false, false)
		ml.add("MASK", int32(255), true, false)
		var buf bytes.Buffer
		if err := ml.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#define FOO bar\n#define SIZE 42\n#define MASK 0x000000ff\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("error_nil_value", func(t *testing.T) {
		var ml macroList
		ml.add("BAD", nil, false, false)
		var buf bytes.Buffer
		err := ml.write(&buf)
		if err == nil {
			t.Fatal("expected error for nil value")
		}
	})

	t.Run("error_invalid_type", func(t *testing.T) {
		var ml macroList
		ml.add("BAD", complex(1, 2), false, false)
		var buf bytes.Buffer
		err := ml.write(&buf)
		if err == nil {
			t.Fatal("expected error for unsupported type")
		}
	})
}
