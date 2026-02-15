package renderer

import "io"

type Renderer interface {
	AddInclude(path string, system bool)
	AddMacro(identifier string, value any, hex bool, raw bool)
	AddData(identifier string, value any, attributes []string, strWidth *int)
	Write(w io.Writer) error
}
