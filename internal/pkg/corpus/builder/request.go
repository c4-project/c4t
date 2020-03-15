// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// Request is the type of requests to a Builder.
type Request struct {
	// Name is the name of the subject to add or modify
	Name string

	// Req is the request payload, which will be one of the *Req structs.
	Req interface{}
}

// SendTo tries to send this request down ch while checking to see if ctx has been cancelled.
func (b Request) SendTo(ctx context.Context, ch chan<- Request) error {
	select {
	case ch <- b:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Add is a request to add the given subject to the corpus.
type Add subject.Subject

// SendAdd tries to send an add request for s down ch, failing if ctx has terminated.
func SendAdd(ctx context.Context, ch chan<- Request, s *subject.Named) error {
	return Request{Name: s.Name, Req: Add(s.Subject)}.SendTo(ctx, ch)
}

// Compile is a request to add the given compiler result to the named subject.
type Compile struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID model.ID

	// Result is the compile result.
	Result subject.CompileResult
}

// Harness is a request to add the given harness to the named subject, under the named architecture.
type Harness struct {
	// Arch is the ID of the architecture for which this lifting is occurring.
	Arch model.ID

	// Harness is the produced harness pathset.
	Harness subject.Harness
}

// Run is a request to add the given run result to the named subject.
type Run struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID model.ID

	// Run is the run result.
	Result subject.Run
}
