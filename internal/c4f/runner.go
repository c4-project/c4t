// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f

import (
	"context"
	"errors"
	"io"

	"github.com/c4-project/c4t/internal/model/service"
)

// ErrNoBaseRunner occurs if we try to use a Runner that has no Runner.Base set.
var ErrNoBaseRunner = errors.New("no base runner supplied")

// Runner stores information about how to run the core ACT binaries.
type Runner struct {
	// DuneExec toggles whether ACT should be run through dune.
	DuneExec bool

	// Base is the basic service runner the ACT runner is using to run binaries.
	Base service.Runner
}

// CmdSpec holds all information about the invocation of an ACT command.
type CmdSpec struct {
	// Cmd is the name of the ACT command (binary) to run.
	Cmd string
	// Subcmd is the name of the ACT subcommand to run.
	Subcmd string
	// Args is the argument vector to supply to the ACT subcommand.
	Args []string
	// Stdout, if given, redirects the command's stdout to this writer.
	Stdout io.Writer
}

// FullArgv gets the full argument vector for the command, including the subcommand.
func (c CmdSpec) FullArgv() []string {
	// Reserving room for the subcommand.
	fargv := make([]string, 1, 1+len(c.Args))
	fargv[0] = c.Subcmd
	return append(fargv, c.Args...)
}

func (a *Runner) Run(ctx context.Context, s CmdSpec) error {
	r := a.Base
	if r == nil {
		return ErrNoBaseRunner
	}
	if s.Stdout != nil {
		r = r.WithStdout(s.Stdout)
	}

	return r.Run(ctx, liftDuneExec(a.DuneExec, s.Cmd, s.FullArgv()))
}

func liftDuneExec(duneExec bool, cmd string, argv []string) service.RunInfo {
	if duneExec {
		cmd, argv = "dune", append([]string{"exec", cmd, "--"}, argv...)
	}
	return *service.NewRunInfo(cmd, argv...)
}
