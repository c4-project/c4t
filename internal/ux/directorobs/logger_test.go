// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/act-tester/internal/ux/directorobs"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// ExampleInstanceLogger_OnArchive is a runnable example for OnArchive.
func ExampleInstanceLogger_OnArchive() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	i, _ := l.Instance(id.FromString("localhost"))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		saver.OnArchiveStart("subj", "subj.tar.gz", 2, i)
		saver.OnArchiveFileAdded("subj", "a.out", 0, i)
		saver.OnArchiveFileMissing("subj", "compile.log", 1, i)
		saver.OnArchiveFinish("subj", i)
		cancel()
	}()
	_ = l.Run(ctx)

	// Output:
	// saving (run [ #0 (Jan  1 00:00:00)]) subj to subj.tar.gz
	// when saving (run [ #0 (Jan  1 00:00:00)]) subj: missing file compile.log
}
