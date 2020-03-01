// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// BuilderReq is the type of requests to a Builder.
type BuilderReq struct {
	// Name is the name of the subject to add or modify
	Name string

	// Req is the request payload, which will be one of the *Req structs.
	Req interface{}
}

// SendTo tries to send this request down ch while checking to see if ctx has been cancelled.
func (b BuilderReq) SendTo(ctx context.Context, ch chan<- BuilderReq) error {
	select {
	case ch <- b:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AddReq is a request to add the given subject to the corpus.
type AddReq subject.Subject

// SendAdd tries to send an add request for s down ch, failing if ctx has terminated.
func SendAdd(ctx context.Context, ch chan<- BuilderReq, s *subject.Named) error {
	return BuilderReq{Name: s.Name, Req: AddReq(s.Subject)}.SendTo(ctx, ch)
}

// AddCompileReq is a request to add the given compiler result to the named subject.
type AddCompileReq struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID model.ID

	// Result is the compile result.
	Result subject.CompileResult
}

// AddHarnessReq is a request to add the given harness to the named subject, under the named architecture.
type AddHarnessReq struct {
	// Arch is the ID of the architecture for which this lifting is occurring.
	Arch model.ID

	// Harness is the produced harness pathset.
	Harness subject.Harness
}

// AddRunReq is a request to add the given run result to the named subject.
type AddRunReq struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID model.ID

	// Run is the run result.
	Result subject.Run
}
