package fuzzer

import (
	"context"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// job contains state for a single fuzzer batch-job.
type job struct {
	// Subject contains the subject for which this job is responsible.
	Subject subject.Subject

	// Driver is the low-level fuzzer.
	Driver SingleFuzzer

	// SubjectCycles is the number of times each subject should be fuzzed.
	SubjectCycles int

	// Pathset points to the pathset to use to work out where to store fuzz output.
	Pathset SubjectPather

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- subject.Subject
}

// FuzzFileset performs a single fuzzing job.
func (j *job) Fuzz(ctx context.Context) error {
	for i := 0; i < j.SubjectCycles; i++ {
		if err := j.fuzzCycle(ctx, i); err != nil {
			return err
		}
	}
	return nil
}

func (j *job) fuzzCycle(ctx context.Context, cycle int) error {
	sc := SubjectCycle{Name: j.Subject.Name, Cycle: cycle}
	spaths := j.Pathset.SubjectPaths(sc)
	if err := j.Driver.FuzzSingle(ctx, j.Rng.Int31(), j.Subject.Litmus, spaths.Litmus, spaths.Trace); err != nil {
		return err
	}
	j.Subject.Fuzz = &spaths
	if err := j.sendSubject(ctx); err != nil {
		return err
	}
	return nil
}

// sendSubject tries to send this job's subject down its result channel.
func (j *job) sendSubject(ctx context.Context) error {
	select {
	case j.ResCh <- j.Subject:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
