// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package query

import (
	"io"
	"io/ioutil"

	"github.com/MattWindsor91/act-tester/internal/controller/query"

	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

func App(outw, errw io.Writer) *c.App {
	return &c.App{
		Name:  "act-tester-query",
		Usage: "performs human-readable queries on a plan file",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
		Writer:                 outw,
		ErrWriter:              errw,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.PlanFileCliFlag(),
		// TODO(@MattWindsor91): template stuff
	}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	pf := stdflag.PlanFileFromCli(ctx)
	q := query.Config{Out: outw}
	return view.RunOnPlanFile(ctx.Context, &q, pf, ioutil.Discard)
}
