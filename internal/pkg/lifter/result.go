// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/cheggaaa/pb/v3"
)

// result is the type of results from the parallelised lifting process.
type result struct {
	// MArch is the machine-qualified architecture for which this lifting is occurring.
	MArch model.MachQualID

	// Subject is the name of the subject that has been lifted.
	Subject string

	// Harness is the produced harness pathset.
	Harness subject.Harness
}

// handleResults waits for nresult results to come in through resCh.
// For each result, it propagates the harness pathset to its subject in c.
func handleResults(ctx context.Context, c corpus.Corpus, nresult int, resCh <-chan result) error {
	bar := pb.StartNew(nresult)
	defer bar.Finish()

	for i := 0; i < nresult; i++ {
		select {
		case r := <-resCh:
			if err := handleResult(r, c); err != nil {
				return err
			}
			bar.Increment()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// handleResult applies a result's liftings to its own subject.
// We make sure to do this sequentially in a single result-handling goroutine, to avoid races.
func handleResult(r result, c corpus.Corpus) error {
	s := c[r.Subject]
	if err := s.AddHarness(r.MArch, r.Harness); err != nil {
		return err
	}
	c[r.Subject] = s
	return nil
}
