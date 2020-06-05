// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"strconv"
	"time"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/timeout"
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
	args := []string{
		"-" + FlagUseJSON,
		"-" + FlagOutDir, m.maybeOverrideDir(dir),
		"-" + FlagCompilerTimeoutLong, m.Config.Quantities.Compiler.Timeout.String(),
		"-" + FlagRunTimeoutLong, m.Config.Quantities.Runner.Timeout.String(),
		"-" + FlagWorkerCountLong, strconv.Itoa(m.Config.Quantities.Runner.NWorkers),
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

// MachConfigCliFlags gets the cli flags for setting up the 'user config' part of a mach or invoker invocation.
func MachConfigCliFlags() []c.Flag {
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
			Name:    FlagCompilerTimeoutLong,
			Aliases: []string{FlagCompilerTimeout},
			Value:   1 * time.Minute,
			Usage:   "a `timeout` to apply to each compilation",
		},
		&c.DurationFlag{
			Name:    FlagRunTimeoutLong,
			Aliases: []string{FlagRunTimeout},
			Value:   1 * time.Minute,
			Usage:   "a `timeout` to apply to each run",
		},
		// TODO(@MattWindsor91): split into compile worker and run worker
		WorkerCountCliFlag(),
		OutDirCliFlag(defaultOutDir),
	}
}

const defaultOutDir = "mach_results"

// MachCliFlags is MachConfigCliFlags plus the dir and JSON flags.
func MachCliFlags() []c.Flag {
	return append(
		MachConfigCliFlags(),
		&c.BoolFlag{
			Name:    FlagUseJSONLong,
			Aliases: []string{FlagUseJSON},
			Usage:   "emit progress reports in JSON form on stderr",
		},
	)
}

// MachConfigFromCli creates a machine configuration using the flags in ctx and the default quantities in defq.
func MachConfigFromCli(ctx *c.Context, defq mach.QuantitySet) mach.UserConfig {
	defq.Override(makeQuantitySet(ctx))

	return mach.UserConfig{
		OutDir:       OutDirFromCli(ctx),
		SkipCompiler: ctx.Bool(FlagSkipCompiler),
		SkipRunner:   ctx.Bool(FlagSkipRunner),
		Quantities:   defq,
	}
}

func makeQuantitySet(ctx *c.Context) mach.QuantitySet {
	return mach.QuantitySet{
		Compiler: compiler.QuantitySet{
			Timeout: timeout.Timeout(ctx.Duration(FlagCompilerTimeoutLong)),
		},
		Runner: runner.QuantitySet{
			Timeout:  timeout.Timeout(ctx.Duration(FlagRunTimeoutLong)),
			NWorkers: ctx.Int(FlagWorkerCountLong),
		},
	}
}
