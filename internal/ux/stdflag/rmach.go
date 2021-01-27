// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"strconv"

	"github.com/c4-project/c4t/internal/quantity"

	c "github.com/urfave/cli/v2"
)

// MachBinName is the name of the machine node binary.
const MachBinName = "c4t-mach"

// MachArgs is the arguments for an invocation of c4t-mach, given directory dir and the config uc.
func MachArgs(dir string, qs quantity.MachNodeSet) []string {
	// We assume that any shell escaping is done elsewhere.
	args := []string{
		"-" + FlagOutDir, dir,
		"-" + FlagCompilerTimeoutLong, qs.Compiler.Timeout.String(),
		"-" + FlagRunTimeoutLong, qs.Runner.Timeout.String(),
		"-" + FlagCompilerWorkerCountLong, strconv.Itoa(qs.Compiler.NWorkers),
		"-" + FlagRunWorkerCountLong, strconv.Itoa(qs.Runner.NWorkers),
	}
	return args
}

// MachInvocation gets the invocation for the local-machine binary as a string list.
func MachInvocation(dir string, qs quantity.MachNodeSet) []string {
	return append([]string{MachBinName}, MachArgs(dir, qs)...)
}

// MachCliFlags gets the cli flags for setting up the 'user config' part of a mach or invoker invocation.
func MachCliFlags() []c.Flag {
	return []c.Flag{
		&c.DurationFlag{
			Name:        FlagCompilerTimeoutLong,
			Aliases:     []string{FlagCompilerTimeout},
			Value:       0,
			Usage:       "a `timeout` to apply to each compilation",
			DefaultText: "from config",
		},
		&c.DurationFlag{
			Name:        FlagRunTimeoutLong,
			Aliases:     []string{FlagRunTimeout},
			Value:       0,
			Usage:       "a `timeout` to apply to each run",
			DefaultText: "from config",
		},
		&c.IntFlag{
			Name:        FlagCompilerWorkerCountLong,
			Aliases:     []string{FlagWorkerCount},
			Value:       0,
			Usage:       "number of compiler `workers` to run in parallel",
			DefaultText: "from config",
		},
		&c.IntFlag{
			Name:        FlagRunWorkerCountLong,
			Aliases:     []string{FlagAltWorkerCount},
			Value:       0,
			Usage:       "number of runner `workers` to run in parallel (not recommended except on manycore machines)",
			DefaultText: "from config",
		},
		OutDirCliFlag(defaultOutDir),
	}
}

const defaultOutDir = "mach_results"

// MachNodeQuantitySetFromCli gets the machine node quantity set from the flags in ctx.
func MachNodeQuantitySetFromCli(ctx *c.Context) quantity.MachNodeSet {
	return quantity.MachNodeSet{
		Compiler: quantity.BatchSet{
			Timeout:  quantity.Timeout(ctx.Duration(FlagCompilerTimeoutLong)),
			NWorkers: ctx.Int(FlagCompilerWorkerCountLong),
		},
		Runner: quantity.BatchSet{
			Timeout:  quantity.Timeout(ctx.Duration(FlagRunTimeoutLong)),
			NWorkers: ctx.Int(FlagRunWorkerCountLong),
		},
	}
}
