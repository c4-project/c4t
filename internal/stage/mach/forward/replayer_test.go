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
	"reflect"
	"strings"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder/mocks"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/forward"
)

// TestReplayer_Run_roundTrip tests a round-trip between Observer and Replayer.
func TestReplayer_Run_roundTrip(t *testing.T) {
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
			recipe.AddFiles("foo.c", "bar.c", "baz.c"),
		))
	com := builder.CompileRequest(
		"foo",
		id.CStyleGCC,
		compilation.CompileResult{
			Result: compilation.Result{Status: status.Ok},
			Files: compilation.CompileFileset{
				Bin: "foo/bin",
				Log: "foo/log",
			},
		})
	run := builder.RunRequest(
		"foo",
		id.CStyleGCC,
		compilation.RunResult{
			Result: compilation.Result{Status: status.Flagged},
		})

	tobs, err := roundTrip(context.Background(), func(obs *forward.Observer) {
		builder.OnBuildStart(m, obs)
		builder.OnBuildRequest(add, obs)
		builder.OnBuildRequest(rec, obs)
		builder.OnBuildRequest(com, obs)
		builder.OnBuildRequest(run, obs)
		builder.OnBuildFinish(obs)
	}, func(obs *mocks.Observer) {
		obs.On("OnBuild", mock.MatchedBy(func(msg builder.Message) bool {
			return msg.Kind == builder.BuildStart &&
				reflect.DeepEqual(*msg.Manifest, m)
		})).Return().Once().On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
			return m.Kind == builder.BuildRequest &&
				m.Request.Name == add.Name &&
				m.Request.Add != nil
		})).Return().Once().On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
			return m.Kind == builder.BuildRequest &&
				m.Request.Name == rec.Name &&
				m.Request.Recipe != nil
		})).Return().Once().On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
			return m.Kind == builder.BuildRequest &&
				m.Request.Name == com.Name &&
				m.Request.Compile != nil
		})).Return().Once().On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
			return m.Kind == builder.BuildRequest &&
				m.Request.Name == run.Name &&
				m.Request.Run != nil
		})).Return().Once().On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
			return m.Kind == builder.BuildFinish
		})).Return().Once()
	})
	require.NoError(t, err)

	tobs.AssertExpectations(t)
}

// TestReplayer_Run_roundTripError tests an error round-trip between Observer and Replayer.
func TestReplayer_Run_roundTripError(t *testing.T) {
	t.Parallel()

	e := fmt.Errorf("it's the end of the world as we know it")

	_, err := roundTrip(context.Background(), func(obs *forward.Observer) {
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

	_, err := roundTrip(ctx, func(*forward.Observer) {}, func(*mocks.Observer) {})
	testhelp.ExpectErrorIs(t, err, ctx.Err(), "replay with immediate cancel")
}

func roundTrip(ctx context.Context, input func(*forward.Observer), obsf func(*mocks.Observer)) (*mocks.Observer, error) {
	pw, obs, tobs, rep := roundTripPipe()
	obsf(tobs)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		input(&obs)
		return pw.Close()
	})
	eg.Go(func() error {
		return rep.Run(ectx)
	})
	return tobs, eg.Wait()
}

func roundTripPipe() (io.Closer, forward.Observer, *mocks.Observer, forward.Replayer) {
	pr, pw := io.Pipe()
	obs := forward.Observer{Encoder: json.NewEncoder(pw)}
	tobs := mocks.Observer{}
	rep := forward.Replayer{Decoder: json.NewDecoder(pr), Observers: []builder.Observer{&tobs}}
	return pw, obs, &tobs, rep
}
