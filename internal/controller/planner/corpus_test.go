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
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/controller/planner"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

type TestProber struct {
	err    error
	probes map[string]struct{}
}

// ProbeSubject is a mock implementation of subject probing.
func (t *TestProber) ProbeSubject(_ context.Context, litmus string) (subject.Named, error) {
	if t.probes == nil {
		t.probes = make(map[string]struct{})
	}
	t.probes[litmus] = struct{}{}

	s, err := subject.New(litmus, subject.WithThreads(2))
	if err != nil {
		return subject.Named{}, err
	}

	return subject.Named{
		Name:    iohelp.ExtlessFile(litmus),
		Subject: *s,
	}, t.err
}

// TestCorpusPlanner_Plan tests the happy path of Plan using a mock SubjectProber.
func TestCorpusPlanner_Plan(t *testing.T) {
	tp := TestProber{}
	p := makeCorpusPlanner(&tp)
	c, err := p.Plan(context.Background())
	if err != nil {
		t.Fatal("unexpected error in Plan:", err)
	}

	if len(c) != p.Size {
		t.Fatalf("corpus size mismatch: got=%d, want=%d", len(c), p.Size)
	}

	for n, s := range c {
		f := n + ".litmus"

		if len(p.Files) <= sort.SearchStrings(p.Files, f) {
			t.Errorf("unexpected corpus subject %q", n)
		}

		if s.OrigLitmus != f {
			t.Errorf("subject %q file mismatch: got=%q, want=%q", n, s.OrigLitmus, f)
		}

		if _, ok := tp.probes[f]; !ok {
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
		// This should enforce a degree of sampling.
		Size: len(in) / 2,
	}
}
