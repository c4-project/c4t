// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// Request is the type of requests to a Builder.
type Request struct {
	// Name is the name of the subject to add or modify
	Name string `json:"name"`

	// Add is populated if this request is an Add.
	Add *Add `json:"add,omitempty"`

	// Compile is populated if this request is a Compile.
	Compile *Compile `json:"compile,omitempty"`

	// Harness is populated if this request is a Harness.
	Harness *Harness `json:"harness,omitempty"`

	// Run is populated if this request is a Run.
	Run *Run `json:"run,omitempty"`
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

// AddRequest constructs an add-subject request for subject s.
func AddRequest(s *subject.Named) Request {
	a := Add(s.Subject)
	return Request{Name: s.Name, Add: &a}
}

// Compile is a request to add the given compiler result to the named subject.
type Compile struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID id.ID

	// Result is the compile result.
	Result subject.CompileResult
}

// CompileRequest constructs an add-compile request for the subject with name sname, compiler ID cid, and result r.
func CompileRequest(sname string, cid id.ID, r subject.CompileResult) Request {
	return Request{Name: sname, Compile: &Compile{CompilerID: cid, Result: r}}
}

// Harness is a request to add the given harness to the named subject, under the named architecture.
type Harness struct {
	// Arch is the ID of the architecture for which this lifting is occurring.
	Arch id.ID

	// Harness is the produced harness pathset.
	Harness subject.Harness
}

// HarnessRequest constructs an add-harness request for the subject with name sname, arch ID arch, and harness h.
func HarnessRequest(sname string, arch id.ID, h subject.Harness) Request {
	return Request{Name: sname, Harness: &Harness{Arch: arch, Harness: h}}
}

// Run is a request to add the given run result to the named subject.
type Run struct {
	// CompilerID is the ID of the compiler that produced this result.
	CompilerID id.ID

	// Run is the run result.
	Result subject.Run
}

// RunRequest constructs an add-run request for the subject with name sname, compiler ID cid, and result r.
func RunRequest(sname string, cid id.ID, r subject.Run) Request {
	return Request{Name: sname, Run: &Run{CompilerID: cid, Result: r}}
}
