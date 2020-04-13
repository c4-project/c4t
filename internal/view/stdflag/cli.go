// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"io"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/config"
	// It's 2020, and tools _still_ can't understand the use of 'v2' unless you do silly hacks like this.
	c "github.com/urfave/cli/v2"
)

// OutDirCliFlag sets up an 'output directory' cli flag.
func OutDirCliFlag(defaultdir string) c.Flag {
	return &c.PathFlag{
		Name:  FlagOutDir,
		Value: defaultdir,
		Usage: usageOutDir,
	}
}

// ActRunnerCliFlags gets the 'cli' flags needed to set up an ACT runner.
func ActRunnerCliFlags() []c.Flag {
	return []c.Flag{
		&c.PathFlag{
			Name:      flagActConfFile,
			Usage:     usageActConfFile,
			TakesFile: true,
		},
		&c.BoolFlag{
			Name:  flagActDuneExec,
			Usage: usageActConfFile,
		},
	}
}

// ActRunnerFromCli makes an ACT runner using the flags previously set up by ActRunnerCliFlags.
func ActRunnerFromCli(ctx *c.Context, errw io.Writer) *act.Runner {
	return &act.Runner{
		DuneExec: ctx.Bool(flagActDuneExec),
		ConfFile: ctx.Path(flagActConfFile),
		Stderr:   errw,
	}
}

// ConfFileCliFlag creates a cli flag for the config file.
func ConfFileCliFlag() c.Flag {
	return &c.PathFlag{
		Name:      flagConfFile,
		Usage:     usageActConfFile,
		TakesFile: true,
	}
}

// ConfFileFromCli sets up a Config using the file flag set up by ConfFileCliFlag.
func ConfFileFromCli(ctx *c.Context) (*config.Config, error) {
	cfile := ctx.Path(flagConfFile)
	return config.Load(cfile)
}

// OutDirFromCli gets the output directory set up by OutDirCliFlag.
func OutDirFromCli(ctx *c.Context) string {
	return ctx.Path(FlagOutDir)
}

// PlanFileCliFlag sets up a standard cli flag for loading a plan file into f.
func PlanFileCliFlag() c.Flag {
	return &c.PathFlag{
		Name:      FlagInputFile,
		TakesFile: true,
		Usage:     usagePlanFile,
	}
}

// PlanFileFromCli retrieves a plan file using the file flag set up by PlanFileCliFlag.
func PlanFileFromCli(ctx *c.Context) string {
	return ctx.Path(FlagInputFile)
}
