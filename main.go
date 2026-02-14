package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"rafaelmartins.com/p/synth-datagen/internal/charts"
	"rafaelmartins.com/p/synth-datagen/internal/codegen"
	"rafaelmartins.com/p/synth-datagen/internal/config"
	"rafaelmartins.com/p/synth-datagen/internal/modules"
	"rafaelmartins.com/p/synth-datagen/internal/renderer"
	"rafaelmartins.com/p/synth-datagen/internal/utils"
	"rafaelmartins.com/p/synth-datagen/internal/version"
)

var (
	oConfig  = flag.String("f", "synth-datagen.yml", "path to configuration file")
	oOutput  = flag.String("o", ".", "path to output directory")
	oCharts  = flag.Bool("c", false, "generate charts and exit")
	oVersion = flag.Bool("v", false, "show version and exit")
)

func check(err any) {
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	flag.Parse()

	if *oVersion {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), version.Version)
		os.Exit(0)
	}

	conf, err := config.New(*oConfig)
	check(err)

	modules.SetGlobalParameters(conf.GlobalParameters)

	for _, out := range conf.Outputs {
		var (
			rndr    renderer.Renderer
			outfile string
		)
		if *oCharts {
			if out.ChartsOutput == "" {
				continue
			}
			outfile = filepath.Join(*oOutput, out.ChartsOutput)
			rndr = charts.New(filepath.Base(out.HeaderOutput))
		} else {
			outfile = filepath.Join(*oOutput, out.HeaderOutput)
			rndr = codegen.NewHeader()
		}

		log.Printf("Generating %q ...", outfile)

		for _, inc := range out.Includes {
			rndr.AddInclude(inc.Path, inc.System)
		}

		for _, mac := range out.Macros {
			rndr.AddMacro(mac.Identifier, mac.Value, mac.Hex, mac.Raw)
		}

		for _, v := range out.Variables {
			rndr.AddData(v.Identifier, v.Value, v.Attributes, v.StringWidth)
		}

		for _, mod := range out.Modules {
			check(modules.Render(rndr, mod.Identifier, mod.Name, mod.Parameters, mod.Selectors))
		}

		check(utils.WriteFile(outfile, rndr))
	}
}
