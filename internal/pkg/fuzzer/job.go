package fuzzer

import (
	"context"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// job contains state for a single fuzzer batch-job.
type job struct {
	// Subject contains the subject for which this job is responsible.
	Subject model.Subject

	// Driver is the low-level fuzzer.
	Driver SingleFuzzer

	// SubjectCycles is the number of times each subject should be fuzzed.
	SubjectCycles int

	// Pathset points to the pathset to use to work out where to store fuzz output.
	Pathset SubjectPather

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- model.Subject
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

// Interface for things that can be used to locate a subject path.
// Mainly exists to permit us to mock out Pathset.
type SubjectPather interface {
	// SubjectPaths gets the output and trace paths for a subject.
	SubjectPaths(name string, cycle int) (outp string, tracep string)
}

func (j *job) fuzzCycle(ctx context.Context, cycle int) error {
	outp, tracep := j.Pathset.SubjectPaths(j.Subject.Name, cycle)
	if err := j.Driver.FuzzSingle(j.Rng.Int31(), j.Subject.Litmus, outp, tracep); err != nil {
		return err
	}
	s2 := j.makeSubject(cycle, outp, tracep)
	if err := j.sendSubject(ctx, s2); err != nil {
		return err
	}
	return nil
}

// makeSubject makes the new, fuzzed subject to send back to the batch fuzzer.
func (j *job) makeSubject(cycle int, outp, tracep string) model.Subject {
	return model.Subject{
		Name:       CycledName(j.Subject.Name, cycle),
		OrigLitmus: j.Subject.Litmus,
		Litmus:     outp,
		TracePath:  tracep,
	}
}

// sendSubject tries to send s down this job's result channel.
func (j *job) sendSubject(ctx context.Context, s model.Subject) error {
	select {
	case j.ResCh <- s:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
