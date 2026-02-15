package stringify

import (
	"testing"
)

func TestLpadding(t *testing.T) {
	tests := []struct {
		name     string
		level    uint8
		expected string
	}{
		{"level_0", 0, ""},
		{"level_1", 1, "    "},
		{"level_2", 2, "        "},
		{"level_3", 3, "            "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lpadding(tt.level)
			if got != tt.expected {
				t.Errorf("lpadding(%d) = %q, want %q", tt.level, got, tt.expected)
			}
		})
	}
}

func TestDumpValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		got := dumpValues([]string{}, 0)
		if got != "{}" {
			t.Errorf("got %q, want %q", got, "{}")
		}
	})

	t.Run("single_value", func(t *testing.T) {
		got := dumpValues([]string{"42"}, 0)
		expected := "{\n    42,\n}"
		if got != expected {
			t.Errorf("got %q, want %q", got, expected)
		}
	})

	t.Run("multiple_fit_one_line", func(t *testing.T) {
		got := dumpValues([]string{"1", "2", "3"}, 0)
		expected := "{\n    1, 2, 3,\n}"
		if got != expected {
			t.Errorf("got %q, want %q", got, expected)
		}
	})

	t.Run("with_nesting_level", func(t *testing.T) {
		got := dumpValues([]string{"1"}, 1)
		expected := "    {\n        1,\n    }"
		if got != expected {
			t.Errorf("got %q, want %q", got, expected)
		}
	})

	t.Run("line_wrapping", func(t *testing.T) {
		// each "0x00000001" is 10 chars + ", " = 12. At level 0 padding is 4.
		// 4 + 8*(10+2) = 100 fits in 100, 9th would be 112 so wraps.
		values := make([]string, 16)
		for i := range values {
			values[i] = "0x00000001"
		}
		got := dumpValues(values, 0)
		if got[0] != '{' {
			t.Error("should start with {")
		}
		lines := 0
		for _, c := range got {
			if c == '\n' {
				lines++
			}
		}
		// opening \n + 2 data lines = 3 newlines
		if lines != 3 {
			t.Errorf("expected 3 newlines for 2-line wrapping, got %d", lines)
		}
	})

	t.Run("single_oversized_value", func(t *testing.T) {
		long := "\"" + string(make([]byte, 120)) + "\""
		got := dumpValues([]string{long}, 0)
		// should still produce valid output with the value on its own line
		expected := "{\n    " + long + ",\n}"
		if got != expected {
			t.Errorf("got %q, want %q", got, expected)
		}
	})
}
