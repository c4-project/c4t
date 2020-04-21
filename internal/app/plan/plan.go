// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package plan contains the app definition for act-tester-plan.
package plan

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/view/singleobs"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
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
	usageMach = "ID of machine to use for this test plan"
)

// App creates the act-tester-plan app.
func App(outw, errw io.Writer) *c.App {
	return &c.App{
		Name:  "act-tester-plan",
		Usage: "runs the planning phase of an ACT test standalone",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
		Writer:                 outw,
		ErrWriter:              errw,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.Int64Flag{
			Name:    flagSeed,
			Aliases: []string{flagSeedShort},
			EnvVars: []string{envSeed},
			Usage:   usageSeed,
			Value:   plan.UseDateSeed,
		},
		&c.StringFlag{
			Name:  stdflag.FlagMachine,
			Usage: usageMach,
		},
		&c.IntFlag{
			Name:    flagCorpusSize,
			Aliases: []string{stdflag.FlagNum},
			Usage:   usageCorpusSize,
		},
	}
	return append(ownFlags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	a := stdflag.ActRunnerFromCli(ctx, errw)

	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	pc, err := makePlanConfig(cfg, errw, a, ctx.Int(flagCorpusSize))
	if err != nil {
		return err
	}

	midstr := ctx.String(stdflag.FlagMachine)
	mach, err := getMachine(cfg, midstr)
	if err != nil {
		return err
	}

	fs := ctx.Args().Slice()
	p, err := pc.Plan(ctx.Context, mach, fs, ctx.Int64(flagSeed))
	if err != nil {
		return err
	}

	return p.Dump(outw)
}

func getMachine(cfg *config.Config, midstr string) (plan.NamedMachine, error) {
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return plan.NamedMachine{}, err
	}

	mach, ok := cfg.Machines[midstr]
	if !ok {
		return plan.NamedMachine{}, fmt.Errorf("no such machine: %s", midstr)
	}
	m := plan.NamedMachine{
		ID:      mid,
		Machine: mach.Machine,
	}
	return m, nil
}

func makePlanConfig(c *config.Config, errw io.Writer, a *act.Runner, cs int) (*planner.Config, error) {
	l := log.New(errw, "", 0)
	cfg := planner.Config{
		CorpusSize: cs,
		Source: planner.Source{
			BProbe:     c,
			CLister:    c,
			CInspector: &compiler.CResolve,
			SProbe:     a,
		},
		Logger:    l,
		Observers: observers(l),
	}
	return &cfg, nil
}

func observers(l *log.Logger) planner.ObserverSet {
	return planner.NewObserverSet(singleobs.Planner(l)...)
}
