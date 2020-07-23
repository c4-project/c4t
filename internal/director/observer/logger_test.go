// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/act-tester/internal/director/observer"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	observer2 "github.com/MattWindsor91/act-tester/internal/stage/analyse/observer"
)

// ExampleInstanceLogger_OnArchive is a runnable example for OnArchive.
func ExampleInstanceLogger_OnArchive() {
	l, _ := observer.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout})
	i, _ := l.Instance(id.FromString("localhost"))

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		observer2.OnArchiveStart("subj", "subj.tar.gz", 2, i)
		observer2.OnArchiveFileAdded("subj", "a.out", 0, i)
		observer2.OnArchiveFileMissing("subj", "compile.log", 1, i)
		observer2.OnArchiveFinish("subj", i)
		cancel()
	}()
	_ = l.Run(ctx, cancel)

	// Output:
	// saving (run [ #0 (Jan  1 00:00:00)]) subj to subj.tar.gz
	// when saving (run [ #0 (Jan  1 00:00:00)]) subj: missing file compile.log
}
