// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package planner_test

import (
	"context"
	"errors"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
)

type TestProber struct {
	err    error
	probes sync.Map
}

// ProbeSubject is a mock implementation of subject probing.
func (t *TestProber) ProbeSubject(_ context.Context, path string) (*subject.Named, error) {
	t.probes.Store(path, struct{}{})

	s, err := subject.New(litmus.New(path, litmus.WithThreads(2)))
	if err != nil {
		return nil, err
	}
	return s.AddName(iohelp.ExtlessFile(path)), t.err
}

// TestCorpusPlanner_Plan tests the happy path of Plan using a mock SubjectProber.
func TestCorpusPlanner_Plan(t *testing.T) {
	tp := TestProber{}
	p := makeCorpusPlanner(&tp)
	c, err := p.Plan(context.Background())
	if err != nil {
		t.Fatal("unexpected error in Plan:", err)
	}

	if len(c) != p.Quantities.CorpusSize {
		t.Fatalf("corpus size mismatch: got=%d, want=%d", len(c), p.Quantities.CorpusSize)
	}

	for n, s := range c {
		f := n + ".litmus"

		if len(p.Files) <= sort.SearchStrings(p.Files, f) {
			t.Errorf("unexpected corpus subject %q", n)
		}

		if s.Source.Path != f {
			t.Errorf("subject %q file mismatch: got=%q, want=%q", n, s.Source.Path, f)
		}

		if _, ok := tp.probes.Load(f); !ok {
			t.Errorf("subject %q not probed by prober", n)
		}
	}
}

// TestCorpusPlanner_Plan_ProbeError tests whether an error set during probing bubbles up properly.
func TestCorpusPlanner_Plan_ProbeError(t *testing.T) {
	tp := TestProber{err: errors.New("polarity of neutron flow reversed")}
	p := makeCorpusPlanner(&tp)
	_, err := p.Plan(context.Background())
	testhelp.ExpectErrorIs(t, err, tp.err, "Plan with error returned by prober")
}

// makeCorpusPlanner builds a test CorpusPlanner using tp as the prober.
func makeCorpusPlanner(tp *TestProber) *planner.CorpusPlanner {
	r := rand.New(rand.NewSource(0))
	in := []string{"foo.litmus", "bar.litmus", "baz.litmus", "foobar.litmus", "foobaz.litmus", "barbaz.litmus"}
	sort.Strings(in)
	return &planner.CorpusPlanner{
		Files:  in,
		Prober: tp,
		Rng:    r,
		Quantities: planner.QuantitySet{
			// This should enforce a degree of sampling.
			CorpusSize: len(in) / 2,
			// This should enforce a degree of parallelism.
			NWorkers: len(in) / 2,
		},
	}
}
