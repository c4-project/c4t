// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package corpus_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
	"golang.org/x/sync/errgroup"
)

type testObserver struct {
	done     bool
	nreqs    int
	adds     map[string]struct{}
	compiles map[string][]model.MachQualID
}

func (t *testObserver) OnStart(nreqs int) {
	t.nreqs = nreqs
}

func (t *testObserver) OnAdd(sname string) {
	if t.adds == nil {
		t.adds = map[string]struct{}{}
	}
	t.adds[sname] = struct{}{}
}

func (t *testObserver) OnCompile(sname string, cid model.MachQualID) {
	if t.compiles == nil {
		t.compiles = map[string][]model.MachQualID{}
	}
	t.compiles[sname] = append(t.compiles[sname], cid)
}

func (t *testObserver) OnFinish() {
	t.done = true
}

// TestBuilder_Run_Adds is a long-form test for exercising a corpus builder on an add run.
func TestBuilder_Run_Adds(t *testing.T) {
	obs := testObserver{}

	adds := []subject.Named{
		{
			Name:    "foo",
			Subject: subject.Subject{Threads: 2, Litmus: "foo.litmus"},
		},
		{
			Name:    "bar",
			Subject: subject.Subject{Threads: 1, Litmus: "foo.litmus"},
		},
		{
			Name:    "baz",
			Subject: subject.Subject{Threads: 4, Litmus: "foo.litmus"},
		},
	}

	c := corpus.BuilderConfig{
		Init:  nil,
		NReqs: len(adds),
		Obs:   &obs,
	}

	b, ch, err := corpus.NewBuilder(c)
	if err != nil {
		t.Fatal("unexpected NewBuilder error:", err)
	}

	var got corpus.Corpus

	eg, ectx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		var rerr error
		got, rerr = b.Run(ectx)
		return rerr
	})
	eg.Go(func() error {
		for i := range adds {
			if err := corpus.SendAdd(ectx, ch, &adds[i]); err != nil {
				return err
			}
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		t.Fatal("unexpected error running builder on Adds:", err)
	}

	checkAddObs(t, obs, c)
	checkAddCorpus(t, adds, got)
}

func checkAddCorpus(t *testing.T, adds []subject.Named, got corpus.Corpus) {
	for _, s := range adds {
		sa, ok := got[s.Name]
		if !ok {
			t.Helper()
			t.Error("misplaced add", s.Name)
		} else if !reflect.DeepEqual(sa, s.Subject) {
			t.Helper()
			t.Errorf("add of %s: got subject %v; want %v", s.Name, sa, s.Subject)
		}
	}
}

func checkAddObs(t *testing.T, obs testObserver, c corpus.BuilderConfig) {
	if obs.nreqs != c.NReqs {
		t.Helper()
		t.Errorf("observer told to expect %d requests; want %d", obs.nreqs, c.NReqs)
	}
	if !obs.done {
		t.Helper()
		t.Error("observer not told the builder was done")
	}
}
