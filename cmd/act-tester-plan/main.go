// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/config"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"
)

const (
	flagSeed        = "seed"
	flagSeedShort   = "s"
	usageSeed       = "`seed` to use for any randomised components of this test plan; -1 uses run time as seed"
	flagCorpusSize  = "corpus-size"
	usageCorpusSize = "`number` of corpus files to select for this test plan;\n" +
		"if positive, the planner will use all viable provided corpus files"
	usageMach = "ID of machine to use for this test plan"
)

func main() {
	app := c.App{
		Name:                   "act-tester-plan",
		Usage:                  "runs the planning phase of an ACT test standalone",
		Flags:                  flags(),
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout, os.Stderr)
		},
	}
	view.LogTopError(app.Run(os.Args))
}

func flags() []c.Flag {
	ownFlags := []c.Flag{
		view.ConfFileCliFlag(),
		&c.Int64Flag{
			Name:    flagSeed,
			Aliases: []string{flagSeedShort},
			Usage:   usageSeed,
			Value:   plan.UseDateSeed,
		},
		&c.StringFlag{
			Name:  view.FlagMachine,
			Usage: usageMach,
		},
		&c.IntFlag{
			Name:    flagCorpusSize,
			Aliases: []string{view.FlagNum},
			Usage:   usageCorpusSize,
		},
	}
	return append(ownFlags, view.ActRunnerCliFlags()...)
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	a := view.ActRunnerFromCli(ctx, errw)

	cfg, err := view.ConfFileFromCli(ctx)
	if err != nil {
		return err
	}

	pc, err := makePlanConfig(cfg, errw, a, ctx.Int(flagCorpusSize))
	if err != nil {
		return err
	}

	midstr := ctx.String(view.FlagMachine)
	mid, mach, err := getMachine(cfg, midstr)
	if err != nil {
		return err
	}

	fs := ctx.Args().Slice()
	p, err := pc.Plan(ctx.Context, mid, mach.Machine, fs, ctx.Int64(flagSeed))
	if err != nil {
		return err
	}

	return p.Dump(outw)
}

func getMachine(cfg *config.Config, midstr string) (id.ID, config.Machine, error) {
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return id.ID{}, config.Machine{}, err
	}

	mach, ok := cfg.Machines[midstr]
	if !ok {
		return id.ID{}, config.Machine{}, fmt.Errorf("no such machine: %s", midstr)
	}
	return mid, mach, nil
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
