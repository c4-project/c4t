// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner_test

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/c4-project/c4t/internal/stage/planner"
	"github.com/c4-project/c4t/internal/subject"
)

type TestProber struct {
	err    error
	probes sync.Map
}

// ProbeSubject is a mock implementation of subject probing.
func (t *TestProber) ProbeSubject(_ context.Context, path string) (*subject.Named, error) {
	t.probes.Store(path, struct{}{})

	l, err := litmus.New(path, litmus.WithThreads(2))
	if err != nil {
		return nil, err
	}
	s, err := subject.New(l)
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
	require.NoError(t, err, "unexpected error in Plan")

	require.Len(t, c, len(p.Files), "corpus size mismatch")

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
func makeCorpusPlanner(tp planner.SubjectProber) *planner.CorpusPlanner {
	in := []string{"foo.litmus", "bar.litmus", "baz.litmus", "foobar.litmus", "foobaz.litmus", "barbaz.litmus"}
	sort.Strings(in)
	return &planner.CorpusPlanner{
		Files:  in,
		Prober: tp,
		Quantities: quantity.PlanSet{
			// This should enforce a degree of parallelism.
			NWorkers: len(in) / 2,
		},
	}
}
