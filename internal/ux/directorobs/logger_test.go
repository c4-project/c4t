// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/ux/directorobs"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// ExampleLogger_Instance_onArchive is a runnable example for Instance that exercises sending archive messages.
func ExampleLogger_Instance_onArchive() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	i, _ := l.Instance(id.FromString("localhost"))

	go func() {
		saver.OnArchiveStart("subj", "subj.tar.gz", 2, i)
		saver.OnArchiveFileAdded("subj", "a.out", 0, i)
		saver.OnArchiveFileMissing("subj", "compile.log", 1, i)
		saver.OnArchiveFinish("subj", i)
		// Important, else the logger will keep waiting for the instance to provide observations.
		i.OnInstanceClose()
	}()
	_ = l.Run(context.Background())

	// Output:
	// saving (cycle [ #0 (Jan  1 00:00:00)]) subj to subj.tar.gz
	// when saving (cycle [ #0 (Jan  1 00:00:00)]) subj: missing file compile.log
}

// TestLogger_Run_empty tests that running a logger with no attached instances works out.
func TestLogger_Run_empty(t *testing.T) {
	t.Parallel()

	l, err := directorobs.NewLogger(iohelp.DiscardCloser(), 0)
	require.NoError(t, err, "logger should construct without errors")
	err = l.Run(context.Background())
	require.NoError(t, err, "no channels = no error")
}

// TestLogger_Run_noMessages tests that running a logger with no messages, works out.
func TestLogger_Run_noMessages(t *testing.T) {
	t.Parallel()

	l, err := directorobs.NewLogger(iohelp.DiscardCloser(), 0)
	require.NoError(t, err, "logger should construct without errors")
	i, err := l.Instance(id.FromString("foo"))
	require.NoError(t, err, "instance should construct without errors")
	go func() {
		i.OnInstanceClose()
	}()
	err = l.Run(context.Background())
	require.NoError(t, err, "should have stopped running with no error")
}
