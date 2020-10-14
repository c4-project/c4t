// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder/mocks"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/subject"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus"
	"golang.org/x/sync/errgroup"
)

// TestBuilder_Run_Adds is a long-form test for exercising a corpus builder on an add run.
func TestBuilder_Run_Adds(t *testing.T) {
	var obs mocks.Observer
	obs.Test(t)

	adds := []subject.Named{
		{
			Name:    "foo",
			Subject: *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(2))),
		},
		{
			Name:    "bar",
			Subject: *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(1))),
		},
		{
			Name:    "baz",
			Subject: *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(4))),
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

	onBuild(&obs, observing.BatchStart, func(n int, s string, request *builder.Request) bool {
		return n == c.NReqs && s == c.Name
	}).Return().Once()
	onBuild(&obs, observing.BatchEnd, func(int, string, *builder.Request) bool {
		return true
	}).Return().Once()

	for i, c := range adds {
		i := i
		c := c
		onBuild(&obs, observing.BatchStep, func(n int, _ string, request *builder.Request) bool {
			return i == n &&
				request != nil &&
				request.Name == c.Name &&
				request.Add != nil &&
				request.Add.Source == c.Source
		}).Return().Once()
	}

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

	checkAddCorpus(t, adds, got)

	obs.AssertExpectations(t)
}

func onBuild(m *mocks.Observer, k observing.BatchKind, f func(int, string, *builder.Request) bool) *mock.Call {
	return m.On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
		return m.Kind == k && f(m.Num, m.Name, m.Request)
	}))
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
	want := builder.AddRequest(
		subject.NewOrPanic(litmus.New("blah", litmus.WithThreads(5))).AddName("foo"))

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
