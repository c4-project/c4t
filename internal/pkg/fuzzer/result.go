package fuzzer

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/cheggaaa/pb/v3"
)

func handleResults(ctx context.Context, fuzzed model.Corpus, resCh <-chan model.Subject) error {
	bar := pb.StartNew(len(fuzzed))
	defer bar.Finish()

	for i := range fuzzed {
		select {
		case fuzzed[i] = <-resCh:
			bar.Increment()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
