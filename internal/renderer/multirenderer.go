package renderer

type multirenderer struct {
	r []Renderer
}

func MultiRenderer(r ...Renderer) *multirenderer {
	return &multirenderer{
		r: r,
	}
}

func (mr *multirenderer) AddInclude(path string, system bool) {
	for _, m := range mr.r {
		m.AddInclude(path, system)
	}
}

func (mr *multirenderer) AddMacro(identifier string, value interface{}, hex bool, raw bool) {
	for _, m := range mr.r {
		m.AddMacro(identifier, value, hex, raw)
	}
}

func (mr *multirenderer) AddData(identifier string, value interface{}, attributes []string) {
	for _, m := range mr.r {
		m.AddData(identifier, value, attributes)
	}
}
