package version

import "runtime/debug"

var (
	Version = "git"
)

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("version: build info not available")
	}

	Version = bi.Main.Version
}
