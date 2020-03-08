// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package resolve

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var (
	// ErrUnknownStyle occurs when we ask the resolver for a compiler style of which it isn't aware.
	ErrUnknownStyle = errors.New("unknown compiler style")

	// CResolve is a pre-populated compiler resolver.
	CResolve = CompilerResolver{Compilers: map[string]Compiler{
		"gcc": GCC{DefaultRun: model.CompilerRunInfo{Cmd: "gcc"}},
	}}
)

type Compiler interface {
	Compile(ctx context.Context, arch model.ID, run *model.CompilerRunInfo, j model.CompileJob, errw io.Writer) error
}

type CompilerResolver struct {
	Compilers map[string]Compiler
}

func (c *CompilerResolver) RunCompiler(ctx context.Context, nc *model.NamedCompiler, j model.CompileJob, errw io.Writer) error {
	sstr := nc.Style.String()
	cp, ok := c.Compilers[sstr]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownStyle, sstr)
	}
	return cp.Compile(ctx, nc.Arch, nc.Run, j, errw)
}
