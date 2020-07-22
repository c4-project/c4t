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

	"github.com/MattWindsor91/act-tester/internal/model/machine"

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
	a := c.App{
		Name:  "act-tester-plan",
		Usage: "runs the planning phase of an ACT test standalone",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
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
		stdflag.WorkerCountCliFlag(),
	}
	return append(ownFlags, stdflag.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	a := stdflag.ActRunnerFromCli(ctx, errw)

	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	pc, err := makePlanConfig(cfg, errw, a, ctx.Int(flagCorpusSize), stdflag.WorkerCountFromCli(ctx))
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

	return p.Write(outw, plan.WriteHuman)
}

func getMachine(cfg *config.Config, midstr string) (machine.Named, error) {
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return machine.Named{}, err
	}

	mach, ok := cfg.Machines[midstr]
	if !ok {
		return machine.Named{}, fmt.Errorf("no such machine: %s", midstr)
	}
	m := machine.Named{
		ID:      mid,
		Machine: mach.Machine,
	}
	return m, nil
}

func makePlanConfig(c *config.Config, errw io.Writer, a planner.SubjectProber, cs, nw int) (*planner.Config, error) {
	l := log.New(errw, "", 0)
	cfg := planner.Config{
		Quantities: planner.QuantitySet{
			CorpusSize: cs,
			NWorkers:   nw,
		},
		Source: planner.Source{
			BProbe:     c,
			CLister:    c.Machines,
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
