package selector

import (
	"errors"
	"fmt"
)

type Selector struct {
	selected []string
}

func New(allowed []string, selected []string) (*Selector, error) {
	if len(allowed) == 0 {
		return nil, errors.New("selector: no selector allowed")
	}
	if len(selected) == 0 {
		return nil, errors.New("selector: no selector selected")
	}

	for _, s := range selected {
		found := false
		for _, a := range allowed {
			if a == s {
				found = true
				break
			}
		}
		if !found {
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
		found := false
		for _, a := range s.selected {
			if x == a {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
