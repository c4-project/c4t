package lifter

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/cheggaaa/pb/v3"
)

// result is the type of results from the parallelised lifting process.
type result struct {
	// Machine is the machine for which this lifting is occurring.
	Machine model.Id

	// Arch is the architecture for which this lifting is occurring.
	Arch model.Id

	// Subject is the subject that has been lifted,
	// passed as a pointer to let the result collector modify it in-place.
	Subject *model.Subject

	// Harness is the produced harness pathset.
	Harness model.Harness
}

// handleResults is a goroutine body that waits for nresult results to come in through resCh.
// For each result, it propagates the harness pathset to its subject.
func handleResults(ctx context.Context, nresult int, resCh <-chan result) error {
	bar := pb.StartNew(nresult)
	defer bar.Finish()

	for i := 0; i < nresult; i++ {
		select {
		case r := <-resCh:
			handleResult(r)
			bar.Increment()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// handleResult applies a result's liftings to its own subject.
// We make sure to do this sequentiall in a single result-handling goroutine, to avoid races.
func handleResult(r result) {
	r.Subject.AddHarness(r.Machine, r.Arch, r.Harness)
}
