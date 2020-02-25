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

// AddReq is a request to add the given subject to the corpus.
type AddReq subject.Subject

// SendAdd tries to send an add request for s down ch, failing if ctx has terminated.
func SendAdd(ctx context.Context, ch chan<- BuilderReq, s *subject.Named) error {
	rq := BuilderReq{
		Name: s.Name,
		Req:  AddReq(s.Subject),
	}
	select {
	case ch <- rq:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// AddCompileReq is a request to add the given compiler result to the named subject.
type AddCompileReq struct {
	// CompilerID is the machine-qualified ID of the compiler that produced this result.
	CompilerID model.MachQualID

	// Result is the compile result.
	Result subject.CompileResult
}
