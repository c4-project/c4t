// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach contains the app definition for act-tester-mach.
package mach

import (
	"io"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/app/invoke"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	bimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	cimpl "github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/mach"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/ux"
	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of this binary.
	// The name of this binary is depended-upon by the invoker, so we store it with the rest of the invocation stuff.
	Name = stdflag.MachBinName

	readme = `
   This part of the tester, also known as the 'machine invoker', runs the parts
   of a testing cycle that are specific to the machine-under-test.

   This command's target audience is a pipe, possibly over SSH, connected to an
   instance of the ` + invoke.Name + ` command.  As such, it doesn't make many
   efforts to be user-friendly, and you probably want to use that command
   instead.
`
)

// App creates the act-tester-mach app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:        Name,
		Usage:       "runs the machine-dependent phase of an ACT test",
		Description: strings.TrimSpace(readme),
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	return stdflag.MachCliFlags()
}

func run(ctx *c.Context, outw, errw io.Writer) error {
	m, err := makeMach(ctx, errw)
	if err != nil {
		return err
	}
	return ux.RunOnCliPlan(ctx, m, outw)
}

func makeMach(ctx *c.Context, errw io.Writer) (*mach.Mach, error) {
	errw = iohelp.EnsureWriter(errw)
	fwd := forward.NewObserver(errw)
	return mach.New(
		&cimpl.CResolve,
		&bimpl.BResolve,
		mach.WithUserConfig(stdflag.MachConfigFromCli(ctx, quantity.MachNodeSet{})),
		mach.ForwardTo(fwd),
	)
}
