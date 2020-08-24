// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/stage/mach"
	c "github.com/urfave/cli/v2"
)

const BinMach = "act-tester-mach"

// MachInvoker tells the various invoker runners how to talk to a mach binary,
// passing through a user config in the form of flags.
type MachInvoker struct {
	Config *mach.UserConfig
}

func (m MachInvoker) MachBin() string {
	return BinMach
}

// MachArgs is the arguments for an invocation of act-tester-mach, given directory dir and the config in this invoker.
func (m MachInvoker) MachArgs(dir string) []string {
	// We assume that any shell escaping is done elsewhere.
	args := []string{
		"-" + FlagOutDir, m.maybeOverrideDir(dir),
		"-" + FlagCompilerTimeoutLong, m.Config.Quantities.Compiler.Timeout.String(),
		"-" + FlagRunTimeoutLong, m.Config.Quantities.Runner.Timeout.String(),
		"-" + FlagCompilerWorkerCountLong, strconv.Itoa(m.Config.Quantities.Compiler.NWorkers),
		"-" + FlagRunWorkerCountLong, strconv.Itoa(m.Config.Quantities.Runner.NWorkers),
	}
	if m.Config.SkipCompiler {
		args = append(args, "-"+FlagSkipCompiler)
	}
	if m.Config.SkipRunner {
		args = append(args, "-"+FlagSkipRunner)
	}
	return args
}

func (m MachInvoker) maybeOverrideDir(dir string) string {
	if ystring.IsBlank(dir) {
		return m.Config.OutDir
	}
	return dir
}

// MachCliFlags gets the cli flags for setting up the 'user config' part of a mach or invoker invocation.
func MachCliFlags() []c.Flag {
	return []c.Flag{
		&c.BoolFlag{
			Name:  FlagSkipCompiler,
			Usage: "if given, skip the compiler",
		},
		&c.BoolFlag{
			Name:  FlagSkipRunner,
			Usage: "if given, skip the runner",
		},
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

// MachConfigFromCli creates a machine configuration using the flags in ctx and the default quantities in defq.
func MachConfigFromCli(ctx *c.Context, defq quantity.MachNodeSet) mach.UserConfig {
	defq.Override(makeQuantitySet(ctx))

	return mach.UserConfig{
		OutDir:       OutDirFromCli(ctx),
		SkipCompiler: ctx.Bool(FlagSkipCompiler),
		SkipRunner:   ctx.Bool(FlagSkipRunner),
		Quantities:   defq,
	}
}

func makeQuantitySet(ctx *c.Context) quantity.MachNodeSet {
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
