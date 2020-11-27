// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package setc contains the app definition for c4t-setc.
package setc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler/optlevel"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	cimpl "github.com/MattWindsor91/c4t/internal/serviceimpl/compiler"

	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/ux"

	"github.com/MattWindsor91/c4t/internal/plan"
	"github.com/MattWindsor91/c4t/internal/ux/stdflag"
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

// ErrCompilerMissing occurs when we can't find the compiler with a given name.
var ErrCompilerMissing = errors.New("compiler not found")

func (c *CompilerSetter) Run(_ context.Context, p *plan.Plan) (*plan.Plan, error) {
	cmp, err := getCompiler(p.Compilers, c.cid)
	if err != nil {
		return nil, err
	}
	if err := c.set(&cmp); err != nil {
		return nil, err
	}
	setCompiler(p.Compilers, c.cid, cmp)
	return p, nil
}

func (c *CompilerSetter) set(cnf *compiler.Configuration) error {
	if err := c.setOpt(cnf); err != nil {
		return err
	}
	return c.setMOpt(cnf)
}

// TODO(@MattWindsor91): move some of this to optlevel?

func (c *CompilerSetter) setOpt(cnf *compiler.Configuration) error {
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

func (c *CompilerSetter) setMOpt(cnf *compiler.Configuration) error {
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

// TODO(@MattWindsor91): move all of these onto a 'config map' type.

func getCompiler(m map[string]compiler.Configuration, id id.ID) (compiler.Configuration, error) {
	cmp, ok := m[id.String()]
	if !ok {
		return compiler.Configuration{}, fmt.Errorf("%w: %s", ErrCompilerMissing, id)
	}
	return cmp, nil
}

func setCompiler(m map[string]compiler.Configuration, id id.ID, c compiler.Configuration) {
	m[id.String()] = c
}
