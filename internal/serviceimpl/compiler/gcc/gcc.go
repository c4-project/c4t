// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// GCC represents GCC-style compilers such as GCC and Clang.
type GCC struct {
	// DefaultRun is the default run information for the particular compiler.
	DefaultRun service.RunInfo
}

// RunCompiler compiles j using a GCC-friendly invocation.
func (g GCC) RunCompiler(ctx context.Context, j job.Compile, errw io.Writer) error {
	run := g.DefaultRun
	if nr := j.CompilerRun(); nr != nil {
		run.Override(*nr)
	}
	cmd := exec.CommandContext(ctx, run.Cmd, Args(run, j)...)
	cmd.Stderr = errw
	return cmd.Run()
}

// Args computes the arguments to pass to GCC for running job j with run info run.
// It does not take j's run info into consideration, and assumes this has already been done.
func Args(run service.RunInfo, j job.Compile) []string {
	args := run.Args
	args = AddStringArg(args, "O", j.SelectedOptName())
	args = AddStringArg(args, "m", j.SelectedMOptName())
	args = append(args, "-o", j.Out)
	args = append(args, j.In...)
	return args
}

// AddStringArg adds the argument '-[k][v]' (note lack of equals sign) to args if v is non-blank; else, returns args.
func AddStringArg(args []string, k, v string) []string {
	if ystring.IsBlank(v) {
		return args
	}
	return append(args, fmt.Sprintf("-%s%s", k, v))
}
