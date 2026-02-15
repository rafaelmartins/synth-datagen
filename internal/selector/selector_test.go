package selector

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("nil_allowed", func(t *testing.T) {
		_, err := New(nil, nil)
		if err == nil {
			t.Fatal("expected error for nil allowed")
		}
		if err.Error() != "selector: no selector allowed" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("empty_allowed", func(t *testing.T) {
		_, err := New([]string{}, nil)
		if err == nil {
			t.Fatal("expected error for empty allowed")
		}
		if err.Error() != "selector: no selector allowed" {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("nil_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if s == nil {
			t.Fatal("expected non-nil selector")
		}
		if len(s.selected) != 0 {
			t.Errorf("expected empty selected, got %v", s.selected)
		}
	})

	t.Run("empty_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if s == nil {
			t.Fatal("expected non-nil selector")
		}
		if len(s.selected) != 0 {
			t.Errorf("expected empty selected, got %v", s.selected)
		}
	})

	t.Run("valid_single_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"b"})
		if err != nil {
			t.Fatal(err)
		}
		if len(s.selected) != 1 {
			t.Fatalf("expected 1 selected, got %d", len(s.selected))
		}
		if s.selected[0] != "b" {
			t.Errorf("expected selected[0]=%q, got %q", "b", s.selected[0])
		}
	})

	t.Run("valid_multiple_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if len(s.selected) != 2 {
			t.Fatalf("expected 2 selected, got %d", len(s.selected))
		}
		if s.selected[0] != "a" || s.selected[1] != "c" {
			t.Errorf("expected selected=[a, c], got %v", s.selected)
		}
	})

	t.Run("valid_all_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "b", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if len(s.selected) != 3 {
			t.Fatalf("expected 3 selected, got %d", len(s.selected))
		}
	})

	t.Run("selected_not_allowed", func(t *testing.T) {
		_, err := New([]string{"a", "b"}, []string{"c"})
		if err == nil {
			t.Fatal("expected error for disallowed selector")
		}
		if err.Error() != `selector: "c" selected but not allowed` {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("one_selected_not_allowed", func(t *testing.T) {
		_, err := New([]string{"a", "b"}, []string{"a", "c"})
		if err == nil {
			t.Fatal("expected error for disallowed selector")
		}
		if err.Error() != `selector: "c" selected but not allowed` {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("empty_string_allowed_and_selected", func(t *testing.T) {
		s, err := New([]string{""}, []string{""})
		if err != nil {
			t.Fatal(err)
		}
		if len(s.selected) != 1 {
			t.Fatalf("expected 1 selected, got %d", len(s.selected))
		}
		if s.selected[0] != "" {
			t.Errorf("expected empty string, got %q", s.selected[0])
		}
	})
}

func TestIsSelected(t *testing.T) {
	t.Run("no_args", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, []string{"a"})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected() {
			t.Error("expected false for no arguments")
		}
	})

	t.Run("single_match", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if !s.IsSelected("a") {
			t.Error("expected true for selected item")
		}
	})

	t.Run("single_no_match", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected("b") {
			t.Error("expected false for unselected item")
		}
	})

	t.Run("multiple_all_match", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "b", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if !s.IsSelected("a", "c") {
			t.Error("expected true when all items are selected")
		}
	})

	t.Run("multiple_partial_match", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c"}, []string{"a", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected("a", "b") {
			t.Error("expected false when not all items are selected")
		}
	})

	t.Run("multiple_no_match", func(t *testing.T) {
		s, err := New([]string{"a", "b", "c", "d"}, []string{"a", "c"})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected("b", "d") {
			t.Error("expected false when no items are selected")
		}
	})

	t.Run("empty_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected("a") {
			t.Error("expected false when nothing is selected")
		}
	})

	t.Run("nil_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected("a") {
			t.Error("expected false when nothing is selected")
		}
	})

	t.Run("item_not_in_allowed_but_query_still_false", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, []string{"a"})
		if err != nil {
			t.Fatal(err)
		}
		// "z" was never in allowed, so it can't be in selected
		if s.IsSelected("z") {
			t.Error("expected false for item never in allowed set")
		}
	})

	t.Run("duplicate_args_all_selected", func(t *testing.T) {
		s, err := New([]string{"a", "b"}, []string{"a"})
		if err != nil {
			t.Fatal(err)
		}
		if !s.IsSelected("a", "a") {
			t.Error("expected true for duplicate selected args")
		}
	})

	t.Run("no_args_empty_selected", func(t *testing.T) {
		s, err := New([]string{"a"}, []string{})
		if err != nil {
			t.Fatal(err)
		}
		if s.IsSelected() {
			t.Error("expected false for no arguments with empty selected")
		}
	})
}
