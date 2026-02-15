package codegen

import (
	"bytes"
	"testing"
)

func TestIncludeAdd(t *testing.T) {
	t.Run("system", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		if len(il) != 1 {
			t.Fatalf("expected 1 include, got %d", len(il))
		}
		if il[0].path != "stdint.h" || !il[0].system {
			t.Errorf("unexpected include: path=%q system=%v", il[0].path, il[0].system)
		}
	})

	t.Run("local", func(t *testing.T) {
		var il includeList
		il.add("config.h", false)
		if len(il) != 1 {
			t.Fatalf("expected 1 include, got %d", len(il))
		}
		if il[0].path != "config.h" || il[0].system {
			t.Errorf("unexpected include: path=%q system=%v", il[0].path, il[0].system)
		}
	})

	t.Run("dedup_same_system", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		il.add("stdint.h", true)
		if len(il) != 1 {
			t.Fatalf("expected 1 include after dedup, got %d", len(il))
		}
		if !il[0].system {
			t.Error("expected system to remain true")
		}
	})

	t.Run("dedup_same_local", func(t *testing.T) {
		var il includeList
		il.add("config.h", false)
		il.add("config.h", false)
		if len(il) != 1 {
			t.Fatalf("expected 1 include after dedup, got %d", len(il))
		}
		if il[0].system {
			t.Error("expected system to remain false")
		}
	})

	t.Run("dedup_system_overridden_by_local", func(t *testing.T) {
		var il includeList
		il.add("foo.h", true)
		il.add("foo.h", false)
		if len(il) != 1 {
			t.Fatalf("expected 1 include after dedup, got %d", len(il))
		}
		if il[0].system {
			t.Error("expected system include to be overridden to local")
		}
	})

	t.Run("dedup_local_not_overridden_by_system", func(t *testing.T) {
		var il includeList
		il.add("foo.h", false)
		il.add("foo.h", true)
		if len(il) != 1 {
			t.Fatalf("expected 1 include after dedup, got %d", len(il))
		}
		if il[0].system {
			t.Error("expected local include to stay local")
		}
	})

	t.Run("different_paths", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		il.add("config.h", false)
		il.add("stdbool.h", true)
		if len(il) != 3 {
			t.Fatalf("expected 3 includes, got %d", len(il))
		}
	})
}

func TestIncludeWrite(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var il includeList
		var buf bytes.Buffer
		if err := il.write(&buf); err != nil {
			t.Fatal(err)
		}
		if buf.String() != "" {
			t.Errorf("expected empty output, got %q", buf.String())
		}
	})

	t.Run("system", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		var buf bytes.Buffer
		if err := il.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#include <stdint.h>\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("local", func(t *testing.T) {
		var il includeList
		il.add("config.h", false)
		var buf bytes.Buffer
		if err := il.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#include \"config.h\"\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("multiple_mixed", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		il.add("config.h", false)
		il.add("stdbool.h", true)
		var buf bytes.Buffer
		if err := il.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#include <stdint.h>\n#include \"config.h\"\n#include <stdbool.h>\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})

	t.Run("dedup_preserves_order", func(t *testing.T) {
		var il includeList
		il.add("stdint.h", true)
		il.add("config.h", false)
		il.add("stdint.h", true) // dedup, no effect
		var buf bytes.Buffer
		if err := il.write(&buf); err != nil {
			t.Fatal(err)
		}
		expected := "\n#include <stdint.h>\n#include \"config.h\"\n"
		if buf.String() != expected {
			t.Errorf("got %q, want %q", buf.String(), expected)
		}
	})
}
