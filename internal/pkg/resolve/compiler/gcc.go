// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/job"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"
)

// GCC represents GCC-style compilers such as GCC and Clang.
type GCC struct {
	// DefaultRun is the default run information for the particular compiler.
	DefaultRun service.RunInfo
}

// RunCompiler compiles j according to run using a GCC-friendly invocation.

func (g GCC) RunCompiler(ctx context.Context, j job.Compile, errw io.Writer) error {
	orun := g.DefaultRun.Override(j.Compiler.Run)
	args := GCCArgs(orun, j)
	cmd := exec.CommandContext(ctx, orun.Cmd, args...)
	cmd.Stderr = errw
	return cmd.Run()
}

// GCCArgs computes the arguments to pass to GCC for running job j with run info run.
func GCCArgs(run service.RunInfo, j job.Compile) []string {
	args := run.Args
	args = append(args, "-o", j.Out)
	args = append(args, j.In...)
	return args
}
