package codegen

import (
	"fmt"
	"io"
)

type include struct {
	path   string
	system bool
}

type includeList []*include

func (i *includeList) add(path string, system bool) {
	for _, inc := range *i {
		if path == inc.path {
			if inc.system && !system {
				inc.system = system
			}
			return
		}
	}

	*i = append(*i, &include{
		path:   path,
		system: system,
	})
}

func (i includeList) write(w io.Writer) error {
	if len(i) > 0 {
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	for _, inc := range i {
		if inc.system {
			if _, err := fmt.Fprintf(w, "#include <%s>\n", inc.path); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, "#include \"%s\"\n", inc.path); err != nil {
				return err
			}
		}
	}

	return nil
}
