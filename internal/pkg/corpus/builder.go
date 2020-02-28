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

	// obs is the observer for the builder.
	obs BuilderObserver

	// reqCh is the receiving channel for requests.
	reqCh <-chan BuilderReq
}

// NewBuilder constructs a Builder according to cfg.
// It fails if the number of target requests is negative.
func NewBuilder(cfg BuilderConfig) (*Builder, chan<- BuilderReq, error) {
	if cfg.NReqs <= 0 {
		return nil, nil, fmt.Errorf("%w: %d", ErrBadBuilderTarget, cfg.NReqs)
	}

	reqCh := make(chan BuilderReq)
	b := Builder{
		corpus: initCorpus(cfg.Init, cfg.NReqs),
		n:      cfg.NReqs,
		obs:    obsOrDefault(cfg.Obs),
		reqCh:  reqCh,
	}
	return &b, reqCh, nil
}

// obsOrDefault fills in a default observer if o is nil.
func obsOrDefault(o BuilderObserver) BuilderObserver {
	if o == nil {
		return SilentObserver{}
	}
	return o
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
	b.obs.OnStart(b.n)
	defer b.obs.OnFinish()

	for i := 0; i < b.n; i++ {
		if err := b.runSingle(ctx); err != nil {
			return nil, err
		}
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
		return b.add(r.Name, subject.Subject(rq))
	case AddCompileReq:
		return b.addCompile(r.Name, rq.CompilerID, rq.Result)
	case AddHarnessReq:
		return b.addHarness(r.Name, rq.Arch, rq.Harness)
	default:
		return fmt.Errorf("%w: %s", ErrBadBuilderRequest, reflect.TypeOf(r.Req).Name())
	}
}

func (b *Builder) add(name string, s subject.Subject) error {
	if err := b.corpus.Add(subject.Named{Name: name, Subject: s}); err != nil {
		return err
	}
	b.obs.OnAdd(name)
	return nil
}

func (b *Builder) addCompile(name string, cid model.ID, res subject.CompileResult) error {
	if err := b.rmwSubject(name, func(s *subject.Subject) error {
		return s.AddCompileResult(cid, res)
	}); err != nil {
		return err
	}
	b.obs.OnCompile(name, cid)
	return nil
}

func (b *Builder) addHarness(name string, arch model.ID, h subject.Harness) error {
	if err := b.rmwSubject(name, func(s *subject.Subject) error {
		return s.AddHarness(arch, h)
	}); err != nil {
		return err
	}
	b.obs.OnHarness(name, arch)
	return nil
}

// rmwSubject hoists a mutating function over subjects so that it operates on the corpus subject name.
// This hoisting function is necessary because we can't directly mutate the subject in-place.
func (b *Builder) rmwSubject(name string, f func(*subject.Subject) error) error {
	s, ok := b.corpus[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrBadBuilderName, name)
	}
	if err := f(&s); err != nil {
		return err
	}
	b.corpus[name] = s
	return nil
}
