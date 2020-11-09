// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package forward_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/quantity"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/subject"

	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer/mocks"
)

// TestReplayer_Run_roundTripBuilder a round-trip between Observer and Replayer over builder requests.
func TestReplayer_Run_roundTripBuilder(t *testing.T) {
	t.Parallel()

	m := builder.Manifest{
		Name:  "test",
		NReqs: 3,
	}

	add := builder.AddRequest(subject.NewOrPanic(litmus.New("foo.litmus")).AddName("foo"))
	rec := builder.RecipeRequest(
		"foo",
		id.ArchX8664,
		recipe.New(
			"recipe",
			recipe.OutNothing,
			recipe.AddFiles("foo.c", "bar.c", "baz.c"),
		))
	com := builder.CompileRequest(
		compilation.Name{SubjectName: "foo", CompilerID: id.CStyleGCC},
		compilation.CompileResult{
			Result: compilation.Result{Status: status.Ok},
			Files: compilation.CompileFileset{
				Bin: "foo/bin",
				Log: "foo/log",
			},
		})
	run := builder.RunRequest(
		compilation.Name{SubjectName: "foo", CompilerID: id.CStyleGCC},
		compilation.RunResult{
			Result: compilation.Result{Status: status.Flagged},
		})

	tobs, err := roundTrip(t, context.Background(), func(obs *forward.Observer) {
		builder.OnBuildStart(m, obs)
		builder.OnBuildRequest(0, add, obs)
		builder.OnBuildRequest(1, rec, obs)
		builder.OnBuildRequest(2, com, obs)
		builder.OnBuildRequest(3, run, obs)
		builder.OnBuildFinish(obs)
	}, func(obs *mocks.Observer) {
		onBuild(obs, observing.BatchStart, func(i int, name string, _ *builder.Request) bool {
			return i == m.NReqs && name == m.Name
		}).Return().Once()
		onBuild(obs, observing.BatchStep, func(_ int, _ string, r *builder.Request) bool {
			return r.Name == add.Name && r.Add != nil
		}).Return().Once()
		onBuild(obs, observing.BatchStep, func(_ int, _ string, r *builder.Request) bool {
			return r.Name == rec.Name && r.Recipe != nil
		}).Return().Once()
		onBuild(obs, observing.BatchStep, func(_ int, _ string, r *builder.Request) bool {
			return r.Name == com.Name && r.Compile != nil
		}).Return().Once()
		onBuild(obs, observing.BatchStep, func(_ int, _ string, r *builder.Request) bool {
			return r.Name == run.Name && r.Run != nil
		}).Return().Once()
		onBuild(obs, observing.BatchEnd, func(_ int, _ string, r *builder.Request) bool {
			return true
		}).Return().Once()
	})
	require.NoError(t, err)

	tobs.AssertExpectations(t)
}

// TestReplayer_Run_roundTripMachNode a round-trip between Observer and Replayer over machine node observations.
func TestReplayer_Run_roundTripMachNode(t *testing.T) {
	t.Parallel()

	qs := quantity.MachNodeSet{
		Compiler: quantity.BatchSet{
			Timeout:  1,
			NWorkers: 2,
		},
		Runner: quantity.BatchSet{
			Timeout:  3,
			NWorkers: 4,
		},
	}

	tobs, err := roundTrip(t, context.Background(), func(obs *forward.Observer) {
		observer.OnCompileStart(qs.Compiler, obs)
		observer.OnRunStart(qs.Runner, obs)
	}, func(obs *mocks.Observer) {
		onMachNode(obs, observer.KindCompileStart, func(q *quantity.MachNodeSet) bool {
			return q.Compiler.NWorkers == qs.Compiler.NWorkers && q.Compiler.Timeout == qs.Compiler.Timeout
		}).Return().Once()
		onMachNode(obs, observer.KindRunStart, func(q *quantity.MachNodeSet) bool {
			return q.Runner.NWorkers == qs.Runner.NWorkers && q.Runner.Timeout == qs.Runner.Timeout
		}).Return().Once()
	})
	require.NoError(t, err)

	tobs.AssertExpectations(t)
}

// TestReplayer_Run_roundTripError tests an error round-trip between Observer and Replayer.
func TestReplayer_Run_roundTripError(t *testing.T) {
	t.Parallel()

	e := fmt.Errorf("it's the end of the world as we know it")

	_, err := roundTrip(t, context.Background(), func(obs *forward.Observer) {
		obs.Error(e)
	}, func(*mocks.Observer) {})

	testhelp.ExpectErrorIs(t, err, forward.ErrRemote, "round-tripping an error")

	if !strings.Contains(err.Error(), e.Error()) {
		t.Fatalf("remote error didn't quote original; orig=%v, remote=%v", e, err)
	}
}

// TestReplayerRun_immediateCancel checks that Run bails out immediately if its context has been cancelled.
func TestReplayer_Run_immediateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := roundTrip(t, ctx, func(*forward.Observer) {}, func(*mocks.Observer) {})
	testhelp.ExpectErrorIs(t, err, ctx.Err(), "replay with immediate cancel")
}

func onBuild(m *mocks.Observer, k observing.BatchKind, f func(int, string, *builder.Request) bool) *mock.Call {
	return m.On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
		return m.Kind == k && f(m.Num, m.Name, m.Request)
	}))
}

func onMachNode(m *mocks.Observer, k observer.MessageKind, f func(set *quantity.MachNodeSet) bool) *mock.Call {
	return m.On("OnMachineNodeAction", mock.MatchedBy(func(m observer.Message) bool {
		return m.Kind == k && f(&m.Quantities)
	}))
}

func roundTrip(t *testing.T, ctx context.Context, input func(*forward.Observer), obsf func(*mocks.Observer)) (*mocks.Observer, error) {
	pw, obs, tobs, rep := roundTripPipe(t)
	obsf(tobs)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		input(obs)
		return pw.Close()
	})
	eg.Go(func() error {
		return rep.Run(ectx)
	})
	return tobs, eg.Wait()
}

func roundTripPipe(t *testing.T) (io.Closer, *forward.Observer, *mocks.Observer, forward.Replayer) {
	pr, pw := io.Pipe()
	obs := forward.NewObserver(pw)
	tobs := mocks.Observer{}
	tobs.Test(t)
	rep := forward.Replayer{Decoder: json.NewDecoder(pr), Observers: []observer.Observer{&tobs}}
	return pw, obs, &tobs, rep
}
