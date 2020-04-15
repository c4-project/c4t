// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/MattWindsor91/act-tester/internal/controller/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/view/singleobs"

	bimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	cimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"

	c "github.com/urfave/cli/v2"

	"github.com/MattWindsor91/act-tester/internal/view"
)

func main() {
	app := c.App{
		Name:                   "act-tester-mach",
		Usage:                  "runs the machine-dependent phase of an ACT test",
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
	return stdflag.MachCliFlags()
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	cfg := makeConfig(ctx, outw, errw)
	pfile := stdflag.PlanFileFromCli(ctx)
	return view.RunOnPlanFile(context.Background(), cfg, pfile, outw)
}

func makeConfig(ctx *c.Context, outw, errw io.Writer) *mach.Config {
	cfg := mach.Config{
		CDriver: &cimpl.CResolve,
		RDriver: &bimpl.BResolve,
		Stdout:  outw,
		User:    stdflag.MachConfigFromCli(ctx, mach.QuantitySet{}),
	}
	setLoggerAndObservers(&cfg, errw, ctx.Bool(stdflag.FlagUseJSONLong))
	return &cfg
}

func setLoggerAndObservers(c *mach.Config, errw io.Writer, jsonStatus bool) {
	errw = iohelp.EnsureWriter(errw)

	if jsonStatus {
		c.Logger = nil
		c.Observers = makeJsonObserver(errw)
		return
	}

	c.Logger = log.New(errw, "[mach] ", log.LstdFlags)
	c.Observers = singleobs.Builder(c.Logger)
}

func makeJsonObserver(errw io.Writer) []builder.Observer {
	return []builder.Observer{&forward.Observer{Encoder: json.NewEncoder(errw)}}
}
