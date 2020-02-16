package compiler

import (
	"context"

	"github.com/cheggaaa/pb/v3"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// result logs information about an attempt to compile a subject with a compiler under test.
type result struct {
	// CompilerID is the machine-qualified ID of the compiler that produced this result.
	CompilerID model.MachQualID

	// Subject is the subject that has been lifted,
	// passed as a pointer to let the result collector modify it in-place.
	Subject *subject.Subject

	subject.CompileResult
}

// handleResults waits for nresult results to come in through resCh.
// For each result, it propagates the compilation results to its subject.
func handleResults(ctx context.Context, nresult int, resCh <-chan result) error {
	bar := pb.StartNew(nresult)
	defer bar.Finish()

	for i := 0; i < nresult; i++ {
		select {
		case r := <-resCh:
			if err := handleResult(r); err != nil {
				return err
			}
			bar.Increment()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// handleResult applies a result's compiler result to its own subject.
// We make sure to do this sequentially in a single result-handling goroutine, to avoid races.
func handleResult(r result) error {
	return r.Subject.AddCompileResult(r.CompilerID, r.CompileResult)
}
