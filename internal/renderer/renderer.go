package renderer

import "io"

type Renderer interface {
	AddInclude(path string, system bool)
	AddMacro(identifier string, value interface{}, hex bool, raw bool)
	AddData(identifier string, value interface{}, attributes []string, strWidth *int)
	Write(w io.Writer) error
}
