package main

import (
	"log"
	"path/filepath"

	"github.com/rafaelmartins/synth-datagen/internal/charts"
	"github.com/rafaelmartins/synth-datagen/internal/cli"
	"github.com/rafaelmartins/synth-datagen/internal/codegen"
	"github.com/rafaelmartins/synth-datagen/internal/config"
	"github.com/rafaelmartins/synth-datagen/internal/modules"
	"github.com/rafaelmartins/synth-datagen/internal/renderer"
	"github.com/rafaelmartins/synth-datagen/internal/utils"
	"github.com/rafaelmartins/synth-datagen/internal/version"
)

var (
	oConfig = &cli.StringOption{
		Name:    'f',
		Default: "synth-datagen.yml",
		Help:    "path to configuration file",
		Metavar: "FILE",
	}
	oOutput = &cli.StringOption{
		Name:    'o',
		Default: ".",
		Help:    "path to output directory",
		Metavar: "DIR",
	}
	oCharts = &cli.BoolOption{
		Name: 'c',
		Help: "generate charts",
	}

	cCli = cli.Cli{
		Help:    "A tool that generates C data headers for synthesizer waveforms and algorithms",
		Version: version.Version,
		Options: []cli.Option{
			oConfig,
			oOutput,
			oCharts,
		},
	}
)

func check(err any) {
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func main() {
	cCli.Parse()

	conf, err := config.New(oConfig.GetValue())
	check(err)

	modules.SetGlobalParameters(conf.GlobalParameters)

	for _, out := range conf.Outputs {
		var (
			rndr    renderer.Renderer
			outfile string
		)
		if oCharts.IsSet() {
			if out.ChartsOutput == "" {
				continue
			}
			outfile = filepath.Join(oOutput.GetValue(), out.ChartsOutput)
			rndr = charts.New(filepath.Base(out.HeaderOutput))
		} else {
			outfile = filepath.Join(oOutput.GetValue(), out.HeaderOutput)
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
