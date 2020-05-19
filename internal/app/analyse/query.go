// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

import (
	"io"
	"io/ioutil"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse"

	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:  "act-tester-analyse",
		Usage: "performs human-readable queries on a plan file",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		// TODO(@MattWindsor91): template stuff
	}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	q := analyse.Config{Out: outw}
	return view.RunOnCliPlan(ctx, &q, ioutil.Discard)
}
