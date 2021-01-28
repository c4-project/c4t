// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"github.com/c4-project/c4t/internal/quantity"
	"github.com/c4-project/c4t/internal/stage/fuzzer"
	c "github.com/urfave/cli/v2"
)

const (
	// FlagCompilerTimeoutLong is a long flag for compiler timeout.
	FlagCompilerTimeoutLong = "compiler-timeout"
	// FlagRunTimeoutLong is a long flag for run timeout.
	FlagRunTimeoutLong = "run-timeout"
	// FlagWorkerCountLong is a long flag for arguments that set a worker count.
	FlagWorkerCountLong = "num-workers"
	// FlagCompilerWorkerCountLong is a long flag for arguments that set a compiler worker count.
	FlagCompilerWorkerCountLong = "num-compiler-workers"
	// FlagRunWorkerCountLong is a long flag for arguments that set a runner worker count.
	FlagRunWorkerCountLong = "num-run-workers"

	// TODO(@MattWindsor91): rename xLong/x to x/xShort.
	flagGlobalTimeout  = "global-timeout"
	usageGlobalTimeout = "`duration` after which experiment will be killed"
)

// RootQuantityFlags sets up the root quantity override flags read by RootQuantitiesFromCli.
func RootQuantityCliFlags() []c.Flag {
	return []c.Flag{
		SubjectFuzzesCliFlag(),
		CorpusSizeCliFlag(),
		GlobalTimeoutCliFlag(),
	}
}

// RootQuantitiesFromCli gets from ctx the root-level quantity overrides specified by the user.
func RootQuantitiesFromCli(ctx *c.Context) quantity.RootSet {
	// TODO(@MattWindsor91): disambiguate the corpus size argument
	return quantity.RootSet{
		GlobalTimeout: GlobalTimeoutFromCli(ctx),
		MachineSet: quantity.MachineSet{
			Fuzz: quantity.FuzzSet{
				CorpusSize:    CorpusSizeFromCli(ctx),
				SubjectCycles: SubjectFuzzesFromCli(ctx),
			},
		},
	}
}

// GlobalTimeoutCliFlag sets up a CLI flag for global timeouts.
func GlobalTimeoutCliFlag() c.Flag {
	return &c.DurationFlag{
		Name:        flagGlobalTimeout,
		Usage:       usageGlobalTimeout,
		DefaultText: "no global timeout",
	}
}

// GlobalTimeoutFromCli retrieves any global timeout specified on the command line in ctx.
func GlobalTimeoutFromCli(ctx *c.Context) quantity.Timeout {
	return quantity.Timeout(ctx.Duration(flagGlobalTimeout))
}

// CorpusSizeCliFlag sets up a 'target corpus size' flag.
func CorpusSizeCliFlag() c.Flag {
	return &c.IntFlag{
		Name:        flagCorpusSize,
		Aliases:     []string{FlagNum},
		Value:       0,
		Usage:       usageCorpusSize,
		DefaultText: "all",
	}
}

// CorpusSizeFromCli retrieves a plan file using the file flag set up by CorpusSizeCliFlag.
func CorpusSizeFromCli(ctx *c.Context) int {
	return ctx.Int(flagCorpusSize)
}

// SubjectFuzzesCliFlag sets up a 'number of cycles' flag.
func SubjectFuzzesCliFlag() c.Flag {
	return &c.IntFlag{Name: flagSubjectCycles, Value: fuzzer.DefaultSubjectFuzzes, Usage: usageSubjectFuzzes}
}

// SubjectFuzzesFromCli retrieves a plan file using the file flag set up by SubjectFuzzesCliFlag.
func SubjectFuzzesFromCli(ctx *c.Context) int {
	return ctx.Int(flagSubjectCycles)
}
