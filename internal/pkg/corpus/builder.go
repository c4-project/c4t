// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
	"github.com/cheggaaa/pb/v3"
)

var (
	// ErrBadBuilderTarget occurs when the target request count for a Builder is non-positive.
	ErrBadBuilderTarget = errors.New("number of builder requests must be positive")

	// ErrBadBuilderName occurs when a builder request specifies a name that isn't in the builder's corpus.
	ErrBadBuilderName = errors.New("requested subject name not in builder")

	// ErrBadBuilderRequest occurs when a builder request has an unknown body type.
	ErrBadBuilderRequest = errors.New("unhandled builder request type")
)

// Builder handles the assembly of corpi from asynchronously-constructed subjects.
type Builder struct {
	// corpus is the corpus being built.
	corpus Corpus

	// n is the number of expected requests for the corpus.
	n int

	// reqCh is the receiving channel for requests.
	reqCh <-chan BuilderReq
}

// NewBuilder constructs a Builder for n incoming requests, copying any bindings in the initial corpus init if non-nil.
// It fails if nreqs is negative.
func NewBuilder(init Corpus, nreqs int) (*Builder, chan<- BuilderReq, error) {
	if nreqs <= 0 {
		return nil, nil, fmt.Errorf("%w: %d", ErrBadBuilderTarget, nreqs)
	}

	reqCh := make(chan BuilderReq)
	b := Builder{
		corpus: initCorpus(init, nreqs),
		n:      nreqs,
		reqCh:  reqCh,
	}
	return &b, reqCh, nil
}

func initCorpus(init Corpus, nreqs int) Corpus {
	if init == nil {
		// The requests are probably all going to be add requests, so it's a good starter capacity.
		return make(Corpus, nreqs)
	}
	return init.Copy()
}

// Run runs the builder in context ctx.
// Run is not thread-safe.
func (b *Builder) Run(ctx context.Context) (Corpus, error) {
	// TODO(@MattWindsor91): decouple this
	bar := pb.StartNew(b.n)
	defer bar.Finish()

	for i := 0; i < b.n; i++ {
		if err := b.runSingle(ctx); err != nil {
			return nil, err
		}
		bar.Increment()
	}

	return b.corpus, nil
}

func (b *Builder) runSingle(ctx context.Context) error {
	select {
	case r := <-b.reqCh:
		return b.runRequest(r)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Builder) runRequest(r BuilderReq) error {
	switch rq := r.Req.(type) {
	case AddReq:
		return b.corpus.Add(subject.Named{Name: r.Name, Subject: subject.Subject(rq)})
	case AddCompileReq:
		return b.addCompile(r.Name, rq.CompilerID, rq.Result)
	default:
		return fmt.Errorf("%w: %s", ErrBadBuilderRequest, reflect.TypeOf(r.Req).Name())
	}
}

func (b *Builder) addCompile(name string, cid model.MachQualID, res subject.CompileResult) error {
	s, ok := b.corpus[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrBadBuilderName, name)
	}
	if err := s.AddCompileResult(cid, res); err != nil {
		return err
	}
	b.corpus[name] = s
	return nil
}
