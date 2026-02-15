package selector

import (
	"errors"
	"fmt"
	"slices"
)

type Selector struct {
	selected []string
}

func New(allowed []string, selected []string) (*Selector, error) {
	if len(allowed) == 0 {
		return nil, errors.New("selector: no selector allowed")
	}

	for _, s := range selected {
		if !slices.Contains(allowed, s) {
			return nil, fmt.Errorf("selector: %q selected but not allowed", s)
		}
	}

	return &Selector{
		selected: selected,
	}, nil
}

func (s *Selector) IsSelected(v ...string) bool {
	if len(v) == 0 {
		return false
	}

	for _, x := range v {
		if !slices.Contains(s.selected, x) {
			return false
		}
	}

	return true
}
