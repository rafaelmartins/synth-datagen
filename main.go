package main

import (
	"log"
	"path/filepath"

	"github.com/rafaelmartins/synth-datagen/internal/charts"
	"github.com/rafaelmartins/synth-datagen/internal/codegen"
	"github.com/rafaelmartins/synth-datagen/internal/config"
	"github.com/rafaelmartins/synth-datagen/internal/modules"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/utils"
)

func check(err any) {
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	conf, err := config.New("synth-datagen.yml")
	check(err)

	modules.SetGlobalParameters(conf.GlobalParameters)

	for hname, out := range conf.Output {
		log.Printf("Generating %q ...", hname)

		hdr := codegen.NewHeader()
		cht := (*charts.Charts)(nil)
		rndr := renderer.Renderer(hdr)
		if out.GraphOutput != "" {
			log.Printf("    With charts: %q", out.GraphOutput)
			cht = charts.New(filepath.Base(hname))
			rndr = renderer.MultiRenderer(hdr, cht)
		}

		for _, inc := range out.Includes {
			rndr.AddInclude(inc.Path, inc.System)
		}

		for _, mac := range out.Macros {
			rndr.AddMacro(mac.Identifier, mac.Value, mac.Hex)
		}

		for _, v := range out.Variables {
			rndr.AddData(v.Identifier, v.Value, v.Attributes)
		}

		for _, mod := range out.Modules {
			check(modules.Render(rndr, mod.Identifier, mod.Name, mod.Parameters, mod.Selectors))
		}

		check(utils.WriteFile(hname, hdr))
		if cht != nil {
			check(utils.WriteFile(out.GraphOutput, cht))
		}
	}
}
