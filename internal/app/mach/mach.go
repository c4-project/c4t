// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mach contains the app definition for c4t-mach.
package mach

import (
	"io"
	"strings"

	"github.com/c4-project/c4t/internal/app/invoke"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	br "github.com/c4-project/c4t/internal/serviceimpl/backend/resolver"
	cimpl "github.com/c4-project/c4t/internal/serviceimpl/compiler"
	"github.com/c4-project/c4t/internal/stage/mach"
	"github.com/c4-project/c4t/internal/stage/mach/forward"
	"github.com/c4-project/c4t/internal/ux"
	"github.com/c4-project/c4t/internal/ux/stdflag"
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

// App creates the c4t-mach app.
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
		&br.Resolve,
		mach.OutputDir(stdflag.OutDirFromCli(ctx)),
		mach.OverrideQuantities(stdflag.MachNodeQuantitySetFromCli(ctx)),
		mach.ForwardTo(fwd),
	)
}
