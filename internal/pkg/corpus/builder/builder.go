// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package builder describes a set of types and methods for building corpi asynchronously.
package builder

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

var (
	// ErrBadTarget occurs when the target request count for a Builder is non-positive.
	ErrBadTarget = errors.New("number of builder requests must be positive")

	// ErrBadBuilderName occurs when a builder request specifies a name that isn't in the builder's corpus.
	ErrBadBuilderName = errors.New("requested subject name not in builder")

	// ErrBadBuilderRequest occurs when a builder request has an unknown body type.
	ErrBadBuilderRequest = errors.New("unhandled builder request type")
)

// Builder handles the assembly of corpi from asynchronously-constructed subjects.
type Builder struct {
	// c is the corpus being built.
	c corpus.Corpus

	// n is the number of expected requests for the corpus.
	n int

	// obs is the observer for the builder.
	obs Observer

	// reqCh is the receiving channel for requests.
	reqCh <-chan Request

	// SendCh is the channel to which requests for the builder should be sent.
	SendCh chan<- Request
}

// NewBuilder constructs a Builder according to cfg.
// It fails if the number of target requests is negative.
func NewBuilder(cfg Config) (*Builder, error) {
	if cfg.NReqs <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrBadTarget, cfg.NReqs)
	}

	reqCh := make(chan Request)
	b := Builder{
		c:      initCorpus(cfg.Init, cfg.NReqs),
		n:      cfg.NReqs,
		obs:    obsOrDefault(cfg.Obs),
		reqCh:  reqCh,
		SendCh: reqCh,
	}
	return &b, nil
}

// obsOrDefault fills in a default observer if o is nil.
func obsOrDefault(o Observer) Observer {
	if o == nil {
		return SilentObserver{}
	}
	return o
}

func initCorpus(init corpus.Corpus, nreqs int) corpus.Corpus {
	if init == nil {
		// The requests are probably all going to be add requests, so it's a good starter capacity.
		return make(corpus.Corpus, nreqs)
	}
	return init.Copy()
}

// Run runs the builder in context ctx.
// Run is not thread-safe.
func (b *Builder) Run(ctx context.Context) (corpus.Corpus, error) {
	b.obs.OnStart(b.n)
	defer b.obs.OnFinish()

	for i := 0; i < b.n; i++ {
		if err := b.runSingle(ctx); err != nil {
			return nil, err
		}
	}

	return b.c, nil
}

func (b *Builder) runSingle(ctx context.Context) error {
	select {
	case r := <-b.reqCh:
		return b.runRequest(r)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Builder) runRequest(r Request) error {
	switch rq := r.Req.(type) {
	case Add:
		return b.add(r.Name, subject.Subject(rq))
	case Compile:
		return b.addCompile(r.Name, rq.CompilerID, rq.Result)
	case Harness:
		return b.addHarness(r.Name, rq.Arch, rq.Harness)
	case Run:
		return b.addRun(r.Name, rq.CompilerID, rq.Result)
	default:
		return fmt.Errorf("%w: %s", ErrBadBuilderRequest, reflect.TypeOf(r.Req).Name())
	}
}

func (b *Builder) add(name string, s subject.Subject) error {
	if err := b.c.Add(subject.Named{Name: name, Subject: s}); err != nil {
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
	b.obs.OnCompile(name, cid, res.Success)
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

func (b *Builder) addRun(name string, cid model.ID, r subject.Run) error {
	if err := b.rmwSubject(name, func(s *subject.Subject) error {
		return s.AddRun(cid, r)
	}); err != nil {
		return err
	}
	b.obs.OnRun(name, cid, r.Status)
	return nil
}

// rmwSubject hoists a mutating function over subjects so that it operates on the corpus subject name.
// This hoisting function is necessary because we can't directly mutate the subject in-place.
func (b *Builder) rmwSubject(name string, f func(*subject.Subject) error) error {
	s, ok := b.c[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrBadBuilderName, name)
	}
	if err := f(&s); err != nil {
		return err
	}
	b.c[name] = s
	return nil
}

// ParBuild runs f in a parallelised manner across the subjects in src.
// It uses the responses from f in a Builder, and returns the resulting corpus.
// Note that src may be different from cfg.Init; this is useful when building a new corpus from scratch.
func ParBuild(ctx context.Context, src corpus.Corpus, cfg Config, f func(context.Context, subject.Named, chan<- Request) error) (corpus.Corpus, error) {
	b, err := NewBuilder(cfg)
	if err != nil {
		return nil, err
	}

	var c corpus.Corpus
	perr := src.Par(ctx, 20,
		func(ctx context.Context, named subject.Named) error {
			return f(ctx, named, b.SendCh)
		},
		func(ctx context.Context) error {
			c, err = b.Run(ctx)
			return err
		})
	return c, perr
}
