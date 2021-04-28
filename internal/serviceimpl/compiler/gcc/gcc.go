// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"context"
	"fmt"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/model/service"
)

// GCC represents GCC-style compilers such as GCC and Clang.
type GCC service.ExtClass

// RunCompiler compiles j using a GCC-friendly invocation.
func (g GCC) RunCompiler(ctx context.Context, j compiler.Job, sr service.Runner) error {
	// TODO(@MattWindsor91): this probably should be done before we get to gcc.
	return sr.Run(ctx, g.makeRunInfo(j))
}

func (g GCC) makeRunInfo(j compiler.Job) service.RunInfo {
	run := g.DefaultRunInfo
	if nr := j.CompilerRun(); nr != nil {
		run.Override(*nr)
	}
	run.AppendArgs(Args(j)...)
	return run
}

// Args computes the arguments to pass to GCC for running job j.
// It does not take j's run info into consideration, and assumes this has already been done.
func Args(j compiler.Job) []string {
	var args []string
	args = AddStringArg(args, "O", j.SelectedOptName())
	args = AddStringArg(args, "m", j.SelectedMOptName())
	args = AddKindArg(args, j.Kind)
	args = append(args, "-o", j.Out)
	args = append(args, j.In...)
	return args
}

// AddKindArg adds to args the appropriate GCC argument for achieving the compile kind mentioned in k.
func AddKindArg(args []string, k compiler.Target) []string {
	switch k {
	case compiler.Obj:
		return append(args, "-c")
	default:
		return args
	}
}

// AddStringArg adds the argument '-[k][v]' (note lack of equals sign) to args if v is non-blank; else, returns args.
func AddStringArg(args []string, k, v string) []string {
	if ystring.IsBlank(v) {
		return args
	}
	return append(args, fmt.Sprintf("-%s%s", k, v))
}

// Probe probes for GCC-style compilers, adding them to target.
func (g GCC) Probe(ctx context.Context, sr service.Runner, classId id.ID, target compiler.ConfigMap) error {
	candidates := service.ExtClass(g).ProbeByVersionCommand(ctx, sr, "--version")
	for k := range candidates {
		// TODO(@MattWindsor91): check version
		cid, c, err := g.expandProbedCommand(classId, g.DefaultRunInfo.NewIfDifferent(k))
		if err != nil {
			return err
		}
		// Need to convert to a string, as we're building a raw config map.
		target[cid.String()] = c
	}
	return nil
}

func (g GCC) expandProbedCommand(classId id.ID, run *service.RunInfo) (id.ID, compiler.Config, error) {
	c := compiler.Config{Style: classId, Run: run}
	cid, err := g.makeID(run)
	return cid, c, err
}

func (g GCC) makeID(run *service.RunInfo) (id.ID, error) {
	if run == nil {
		return id.TryFromString(g.DefaultRunInfo.Cmd)
	}
	return run.SystematicID()
}
