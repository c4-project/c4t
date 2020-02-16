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

// Fuzz performs a single fuzzing job.
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
	if err := j.Driver.FuzzSingle(j.Rng.Int31(), j.Subject.Litmus, spaths.FileLitmus, spaths.FileTrace); err != nil {
		return err
	}
	s2 := j.makeSubject(sc, spaths)
	if err := j.sendSubject(ctx, s2); err != nil {
		return err
	}
	return nil
}

// makeSubject makes the new, fuzzed subject to send back to the batch fuzzer.
func (j *job) makeSubject(sc SubjectCycle, ps SubjectPathset) subject.Subject {
	return subject.Subject{
		Name:       sc.String(),
		OrigLitmus: j.Subject.Litmus,
		Litmus:     ps.FileLitmus,
		TracePath:  ps.FileTrace,
	}
}

// sendSubject tries to send s down this job's result channel.
func (j *job) sendSubject(ctx context.Context, s subject.Subject) error {
	select {
	case j.ResCh <- s:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
