// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"context"
	"io"
	"os/exec"
)

// Runner stores information about how to run the core ACT binaries.
type Runner struct {
	// TODO(@MattWindsor91): consider turning DuneExec into a CmdRunner.

	// DuneExec toggles whether ACT should be run through dune.
	DuneExec bool
	// Stderr is the destination for any error output from ACT commands.
	Stderr io.Writer
	// CmdRunner lets one mock out the low-level ACT command runner.
	CmdRunner CmdRunner
}

// CmdRunner is the type of low-level ACT command runners.
type CmdRunner interface {
	// TODO(@MattWindsor91): harmonise with srvrun.ExecRunner?  The two seem to have slightly different purposes.

	// Run runs s; the command will terminate if ctx is cancelled.
	Run(ctx context.Context, s CmdSpec) error
}

//go:generate mockery --name=CmdRunner

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
	// Mocking opportunity.
	if a.CmdRunner != nil {
		return a.CmdRunner.Run(ctx, s)
	}
	return a.runInner(ctx, s)
}

func (a *Runner) runInner(ctx context.Context, s CmdSpec) error {
	fargv := s.FullArgv()
	dcmd, dargv := liftDuneExec(a.DuneExec, s.Cmd, fargv)
	c := exec.CommandContext(ctx, dcmd, dargv...)
	c.Stderr = a.Stderr
	c.Stdout = s.Stdout
	return c.Run()
}

func liftDuneExec(duneExec bool, cmd string, argv []string) (string, []string) {
	if duneExec {
		return "dune", append([]string{"exec", cmd, "--"}, argv...)
	}
	return cmd, argv
}
