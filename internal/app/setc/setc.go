// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package setc contains the app definition for c4t-setc.
package setc

// TODO(@MattWindsor91): make this work more orthogonally with perturber, or merge the two.

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/1set/gut/ystring"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	cimpl "github.com/c4-project/c4t/internal/serviceimpl/compiler"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/ux"

	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	name = "c4t-setc"

	usage = "sets various compiler properties"

	readme = `
   Sets various properties of a compiler in an existing plan file.
`

	flagCompilerLong  = "compiler"
	flagCompilerShort = "c"
	usageCompiler     = "modify the compiler with this `ID`"

	flagOptLong  = "opt-level"
	flagOptShort = "O"
	usageOpt     = "set the compiler's optimisation level `name`"

	flagMoptLong  = "machine-opt"
	flagMoptShort = "m"
	usageMopt     = "set the compiler's machine optimising profile `name`"
)

// App creates the c4t-plan app.
func App(outw, errw io.Writer) *c.App {
	a := c.App{
		Name:        name,
		Usage:       usage,
		Description: readme,
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, os.Stdout)
		},
	}
	return stdflag.SetCommonAppSettings(&a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		&c.StringFlag{
			Name:     flagCompilerLong,
			Aliases:  []string{flagCompilerShort},
			Usage:    usageCompiler,
			Required: true,
		},
		&c.StringFlag{
			Name:    flagOptLong,
			Aliases: []string{flagOptShort},
			Usage:   usageOpt,
		},
		&c.StringFlag{
			Name:    flagMoptLong,
			Aliases: []string{flagMoptShort},
			Usage:   usageMopt,
		},
	}
}

func run(ctx *c.Context, outw io.Writer) error {
	cid, err := compilerID(ctx)
	if err != nil {
		return err
	}

	cs := CompilerSetter{
		inspector: &cimpl.CResolve,
		cid:       cid,
		opt:       ctx.String(flagOptLong),
		mopt:      ctx.String(flagMoptLong),
	}
	return ux.RunOnCliPlan(ctx, &cs, outw)
}

func compilerID(ctx *c.Context) (id.ID, error) {
	cidstr := ctx.String(flagCompilerLong)
	cid, err := id.TryFromString(cidstr)
	if err != nil {
		return id.ID{}, err
	}
	return cid, nil
}

type CompilerSetter struct {
	inspector compiler.Inspector
	cid       id.ID
	opt       string
	mopt      string
}

// Stage gets the stage record for a compiler setter.
func (*CompilerSetter) Stage() stage.Stage {
	return stage.SetCompiler
}

// Close does nothing.
func (*CompilerSetter) Close() error {
	return nil
}

func (c *CompilerSetter) Run(_ context.Context, p *plan.Plan) (*plan.Plan, error) {
	err := p.Compilers.Map(c.cid, func(inst *compiler.Instance) error {
		return c.set(inst)
	})
	return p, err
}

func (c *CompilerSetter) set(cnf *compiler.Instance) error {
	// TODO(@MattWindsor91): allow overriding this.
	cnf.ConfigTime = time.Now()
	// TODO(@MattWindsor91): copy mutant ID?

	if err := c.setOpt(cnf); err != nil {
		return err
	}
	return c.setMOpt(cnf)
}

// TODO(@MattWindsor91): move some of this to optlevel?

func (c *CompilerSetter) setOpt(cnf *compiler.Instance) error {
	if ystring.IsBlank(c.opt) {
		cnf.SelectedOpt = nil
		return nil
	}

	opts, err := compiler.SelectLevels(c.inspector, &cnf.Compiler)
	if err != nil {
		return err
	}
	opt, ok := opts[c.opt]
	if !ok {
		return fmt.Errorf("unknown optimisation level: %s", c.opt)
	}
	cnf.SelectedOpt = &optlevel.Named{Name: c.opt, Level: opt}
	return nil
}

func (c *CompilerSetter) setMOpt(cnf *compiler.Instance) error {
	if ystring.IsBlank(c.mopt) {
		cnf.SelectedMOpt = ""
		return nil
	}

	mopts, err := compiler.SelectMOpts(c.inspector, &cnf.Compiler)
	if err != nil {
		return err
	}
	_, ok := mopts[c.mopt]
	if !ok {
		return fmt.Errorf("unknown machine profile: %s", c.mopt)
	}
	cnf.SelectedMOpt = c.mopt
	return nil
}
