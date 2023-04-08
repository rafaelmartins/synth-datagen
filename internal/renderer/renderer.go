package renderer

type Renderer interface {
	AddInclude(path string, system bool)
	AddMacro(identifier string, value interface{}, hex bool, raw bool)
	AddData(identifier string, value interface{}, attributes []string)
}
