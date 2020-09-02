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

	"github.com/MattWindsor91/act-tester/internal/quantity"

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
	envSeed       = "ACT_SEED"
	flagSeed      = "seed"
	flagSeedShort = "s"
	usageSeed     = "`seed` to use for any randomised components of this test plan"
)

// App creates the act-tester-perturb app.
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
		stdflag.VerboseFlag(),
		stdflag.ConfFileCliFlag(),
		&c.Int64Flag{
			Name:        flagSeed,
			Aliases:     []string{flagSeedShort},
			EnvVars:     []string{envSeed},
			Usage:       usageSeed,
			Value:       plan.UseDateSeed,
			DefaultText: "set seed from time",
		},
		stdflag.CorpusSizeCliFlag(),
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

	qs := quantities(ctx, cfg)

	l := log.New(errw, "[perturb] ", log.LstdFlags)

	return perturber.New(
		&compiler.CResolve,
		perturber.ObserveWith(singleobs.Perturber(l, stdflag.Verbose(ctx))...),
		perturber.OverrideQuantities(qs),
		perturber.UseSeed(ctx.Int64(flagSeed)),
	)
}

func quantities(ctx *c.Context, cfg *config.Config) quantity.PerturbSet {
	qs := cfg.Quantities.Perturb
	qs.Override(quantity.PerturbSet{
		CorpusSize: stdflag.CorpusSizeFromCli(ctx),
	})
	return qs
}
