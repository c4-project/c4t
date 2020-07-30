// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package perturb contains the app definition for act-tester-perturb.
package perturb

import (
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/stage/perturber"
	"github.com/MattWindsor91/act-tester/internal/ux"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/ux/singleobs"
	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	envSeed         = "ACT_SEED"
	flagSeed        = "seed"
	flagSeedShort   = "s"
	usageSeed       = "`seed` to use for any randomised components of this test plan; -1 uses run time as seed"
	flagCorpusSize  = "corpus-size"
	usageCorpusSize = "`number` of corpus files to select for this test plan;\n" +
		"if positive, the planner will use all viable provided corpus files"
)

// App creates the act-tester-plan app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:  "act-tester-perturb",
		Usage: "perturbs a test plan",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.Int64Flag{
			Name:    flagSeed,
			Aliases: []string{flagSeedShort},
			EnvVars: []string{envSeed},
			Usage:   usageSeed,
			Value:   plan.UseDateSeed,
		},
		&c.IntFlag{
			Name:    flagCorpusSize,
			Aliases: []string{stdflag.FlagNum},
			Usage:   usageCorpusSize,
		},
		stdflag.WorkerCountCliFlag(),
	}
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	pr, err := makePerturber(ctx, errw)
	if err != nil {
		return err
	}
	return ux.RunOnCliPlan(ctx, pr, outw)
}

func makePerturber(ctx *c.Context, errw io.Writer) (*perturber.Perturber, error) {
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return nil, err
	}

	qs := quantities(ctx)
	src := source(cfg)

	l := log.New(errw, "[perturb] ", log.LstdFlags)

	return perturber.New(
		src,
		perturber.ObserveWith(singleobs.Perturber(l)...),
		perturber.OverrideQuantities(qs),
		perturber.UseSeed(ctx.Int64(flagSeed)),
	)
}

func source(cfg *config.Config) perturber.Source {
	return perturber.Source{
		CLister:    cfg.Machines,
		CInspector: &compiler.CResolve,
	}
}

func quantities(ctx *c.Context) perturber.QuantitySet {
	return perturber.QuantitySet{
		CorpusSize: ctx.Int(flagCorpusSize),
		NWorkers:   stdflag.WorkerCountFromCli(ctx),
	}
}
