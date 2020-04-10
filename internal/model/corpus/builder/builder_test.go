// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"golang.org/x/sync/errgroup"
)

// TestBuilder_Run_Adds is a long-form test for exercising a corpus builder on an add run.
func TestBuilder_Run_Adds(t *testing.T) {
	obs := builder.MockObserver{}

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

	c := builder.Config{
		Init:      nil,
		Observers: []builder.Observer{&obs},
		Manifest: builder.Manifest{
			Name:  "foobar",
			NReqs: len(adds),
		},
	}

	b, err := builder.New(c)
	if err != nil {
		t.Fatal("unexpected New error:", err)
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
			if err := builder.AddRequest(&adds[i]).SendTo(ectx, b.SendCh); err != nil {
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

func checkAddObs(t *testing.T, obs builder.MockObserver, c builder.Config) {
	if obs.Manifest.NReqs != c.NReqs {
		t.Helper()
		t.Errorf("observer told to expect %d requests; want %d", obs.Manifest.NReqs, c.NReqs)
	}
	if !obs.Done {
		t.Helper()
		t.Error("observer not told the builder was done")
	}
}

func TestBuilderReq_SendTo(t *testing.T) {
	ch := make(chan builder.Request)

	t.Run("success", func(t *testing.T) {
		eg, ectx := errgroup.WithContext(context.Background())
		if err := exerciseSendTo(t, eg, ectx, ch); err != nil {
			t.Error("unexpected error:", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		ctx, c := context.WithCancel(context.Background())
		eg, ectx := errgroup.WithContext(ctx)
		c()
		err := exerciseSendTo(t, eg, ectx, ch)
		testhelp.ExpectErrorIs(t, err, ectx.Err(), "on cancelled SendTo")
	})
}

func exerciseSendTo(t *testing.T, eg *errgroup.Group, ectx context.Context, ch chan builder.Request) error {
	want := builder.AddRequest(&subject.Named{
		Name: "foo",
		Subject: subject.Subject{
			Threads: 5,
			Litmus:  "blah",
		}})

	eg.Go(func() error {
		select {
		case got := <-ch:
			if !reflect.DeepEqual(got, want) {
				t.Helper()
				t.Errorf("received request malformed: got=%v, want=%v", got, want)
			}
			return nil
		case <-ectx.Done():
			return ectx.Err()
		}
	})
	eg.Go(func() error {
		return want.SendTo(ectx, ch)
	})
	return eg.Wait()
}
