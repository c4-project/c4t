// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/subject/corpus"
	"github.com/MattWindsor91/c4t/internal/subject/status"

	"github.com/MattWindsor91/c4t/internal/ux/directorobs"
	"github.com/stretchr/testify/require"
)

// TestNewStatPersister tests NewStatPersister, as well as various other statset manipulations.
func TestNewStatPersister(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	path := filepath.Join(td, "stats.json")
	f, err := directorobs.OpenStatFile(path)
	require.NoError(t, err, "should be able to open stat file in temp dir")

	sp, err := directorobs.NewStatPersister(f)
	require.NoError(t, err, "should be able to open stats persister on new file")

	mid := id.FromString("foo.bar")

	cyc := director.Cycle{MachineID: mid, Iter: 42, Start: time.Now()}
	sp.OnCycle(director.CycleStartMessage(cyc))

	ana := analysis.Analysis{
		ByStatus: map[status.Status]corpus.Corpus{
			status.Ok:             corpus.New("foo", "bar", "baz"),
			status.CompileTimeout: corpus.New("foo", "bar"),
			status.Flagged:        corpus.New("foobaz", "baz"),
		},
		Flags: status.FlagCompileTimeout | status.FlagFlagged,
	}
	sp.OnCycleAnalysis(director.CycleAnalysis{Cycle: cyc, Analysis: ana})

	require.NoError(t, sp.Close(), "should be able to close file")
	// assuming f has been closed here
	f, err = directorobs.OpenStatFile(path)
	require.NoError(t, err, "should be able to reopen stat file")

	var s directorobs.Statset
	if assert.NoError(t, s.Load(f), "should be able to load stats from file") {
		if assert.NotNil(t, s.Machines, "machines should be populated") {
			// Can't compare cycles directly because the time might have been clipped
			assert.Equal(t, cyc.MachineID, s.Machines[mid.String()].LastCycle.MachineID, "last cycle should be the one sent")
			assert.Equal(t, cyc.Iter, s.Machines[mid.String()].LastCycle.Iter, "last cycle should be the one sent")
			for i := status.Ok; i <= status.Last; i++ {
				want := uint64(len(ana.ByStatus[i]))
				assert.Equal(t, want, s.Machines[mid.String()].StatusTotals[i], "status total didn't match")
				assert.Equal(t, want, s.Machines[mid.String()].SessionStatusTotals[i], "session status total didn't match")
			}
		}
	}
	require.NoError(t, f.Close(), "should be able to re-close file")

	// TODO(@MattWindsor91): test session erasure?
}
