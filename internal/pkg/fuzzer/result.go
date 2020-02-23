package fuzzer

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/cheggaaa/pb/v3"
)

func handleResults(ctx context.Context, fuzzed subject.Corpus, nfuzzes int, resCh <-chan subject.Named) error {
	bar := pb.StartNew(nfuzzes)
	defer bar.Finish()

	for i := 0; i < nfuzzes; i++ {
		select {
		case r := <-resCh:
			if err := fuzzed.Add(r); err != nil {
				return err
			}
			bar.Increment()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
