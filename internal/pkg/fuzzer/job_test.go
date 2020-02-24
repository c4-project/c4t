package fuzzer_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
	"golang.org/x/sync/errgroup"
)

// TestJob_Fuzz tests various aspects of a job fuzz.
func TestJob_Fuzz(t *testing.T) {
	resCh := make(chan corpus.BuilderReq)

	j := fuzzer.Job{
		Subject:       subject.Named{Name: "foo"},
		Driver:        fuzzer.NopFuzzer{},
		SubjectCycles: 10,
		Pathset:       fuzzer.NewPathset("test"),
		Rng:           rand.New(rand.NewSource(0)),
		ResCh:         resCh,
	}

	eg, ectx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return j.Fuzz(ectx)
	})
	eg.Go(func() error {
		for i := 0; i < 10; i++ {
			select {
			case r := <-resCh:
				// TODO(@MattWindsor91): other checks
				wname := fuzzer.SubjectCycle{Name: "foo", Cycle: i}.String()
				if r.Name != wname {
					t.Errorf("wrong fuzz result name: got=%q, want=%q", r.Name, wname)
				}
			case <-ectx.Done():
				return ectx.Err()
			}
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		t.Fatalf("unexpected errgroup error: %v", err)
	}
}
