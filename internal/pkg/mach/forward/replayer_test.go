// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package forward_test

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/mach/forward"
)

// TestReplayer_Run_roundTrip tests a round-trip between Observer and Replayer.
func TestReplayer_Run_roundTrip(t *testing.T) {
	t.Parallel()

	m := builder.Manifest{
		Name:  "test",
		NReqs: 3,
	}

	add := builder.AddRequest(
		&subject.Named{
			Name:    "foo",
			Subject: subject.Subject{Litmus: "foo.litmus"},
		})
	harness := builder.HarnessRequest(
		"foo",
		model.ArchX8664,
		subject.Harness{
			Dir:   "harness",
			Files: []string{"foo.c", "bar.c", "baz.c"},
		})
	compile := builder.CompileRequest(
		"foo",
		model.IDFromString("gcc"),
		subject.CompileResult{
			Success: true,
			Files: subject.CompileFileset{
				Bin: "foo/bin",
				Log: "foo/log",
			},
		})
	run := builder.RunRequest(
		"foo",
		model.IDFromString("gcc"),
		subject.Run{
			Status: subject.StatusFlagged,
		})

	tobs, err := roundTrip(context.Background(), func(obs *forward.Observer) {
		obs.OnStart(m)
		obs.OnRequest(add)
		obs.OnRequest(harness)
		obs.OnRequest(compile)
		obs.OnRequest(run)
		obs.OnFinish()
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if !reflect.DeepEqual(tobs.Manifest, m) {
		t.Errorf("manifest mismatch: recv=%v, send=%v", tobs.Manifest, m)
	}
	if !tobs.Done {
		t.Error("test observer didn't receive OnFinish")
	}

	if _, addOk := tobs.Adds[add.Name]; !addOk {
		t.Error("add not propagated")
	}
	if len(tobs.Harnesses[harness.Name]) != 1 {
		t.Error("harness not propagated")
	}
	if len(tobs.Compiles[compile.Name]) != 1 {
		t.Error("compile not propagated")
	}
	if len(tobs.Runs[run.Name]) != 1 {
		t.Error("run not propagated")
	}
}

// TestReplayerRun_immediateCancel checks that Run bails out immediately if its context has been cancelled.
func TestReplayer_Run_immediateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := roundTrip(ctx, func(*forward.Observer) {})
	testhelp.ExpectErrorIs(t, err, ctx.Err(), "replay with immediate cancel")
}

func roundTrip(ctx context.Context, input func(*forward.Observer)) (*builder.MockObserver, error) {
	pw, obs, tobs, rep := roundTripPipe()

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

func roundTripPipe() (io.Closer, forward.Observer, *builder.MockObserver, forward.Replayer) {
	pr, pw := io.Pipe()
	obs := forward.Observer{Encoder: json.NewEncoder(pw)}
	tobs := builder.MockObserver{}
	rep := forward.Replayer{Decoder: json.NewDecoder(pr), Observer: &tobs}
	return pw, obs, &tobs, rep
}
