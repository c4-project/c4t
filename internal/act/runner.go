// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/helper/srvrun"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Runner stores information about how to run the core ACT binaries.
type Runner struct {
	// TODO(@MattWindsor91): consider turning DuneExec into a CmdRunner.

	// DuneExec toggles whether ACT should be run through dune.
	DuneExec bool
	// Stderr is the destination for any error output from ACT commands.
	Stderr io.Writer
	// RunnerFactory lets one mock out the low-level ACT command runner.
	RunnerFactory func(outw, errw io.Writer) service.Runner
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

func execRunner(outw, errw io.Writer) service.Runner {
	return srvrun.NewExecRunner(srvrun.StdoutTo(outw), srvrun.StderrTo(errw))
}

func (a *Runner) Run(ctx context.Context, s CmdSpec) error {
	rf := execRunner
	// Mocking opportunity.
	if a.RunnerFactory != nil {
		rf = a.RunnerFactory
	}
	return a.runInner(ctx, s, rf)
}

func (a *Runner) runInner(ctx context.Context, s CmdSpec, rf func(io.Writer, io.Writer) service.Runner) error {
	fargv := s.FullArgv()
	ri := liftDuneExec(a.DuneExec, s.Cmd, fargv)
	return rf(a.Stderr, s.Stdout).Run(ctx, ri)
}

func liftDuneExec(duneExec bool, cmd string, argv []string) service.RunInfo {
	if duneExec {
		cmd, argv = "dune", append([]string{"exec", cmd, "--"}, argv...)
	}
	return *service.NewRunInfo(cmd, argv...)
}
