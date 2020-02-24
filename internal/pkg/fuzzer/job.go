package fuzzer

import (
	"context"
	"errors"
	"math/rand"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// Job contains state for a single fuzzer batch-Job.
type Job struct {
	// Subject contains the subject for which this Job is responsible.
	Subject subject.Named

	// Driver is the low-level fuzzer.
	Driver SingleFuzzer

	// SubjectCycles is the number of times each subject should be fuzzed.
	SubjectCycles int

	// Pathset points to the pathset to use to work out where to store fuzz output.
	Pathset SubjectPather

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- corpus.BuilderReq
}

// Fuzz performs a single fuzzing Job.
func (j *Job) Fuzz(ctx context.Context) error {
	if err := j.check(); err != nil {
		return err
	}

	for i := 0; i < j.SubjectCycles; i++ {
		if err := j.fuzzCycle(ctx, i); err != nil {
			return err
		}
	}
	return nil
}

// check checks the health of the job before running it.
func (j *Job) check() error {
	if j.Pathset == nil {
		return iohelp.ErrPathsetNil
	}
	if j.Rng == nil {
		return errors.New("RNG nil")
	}
	return nil
}

func (j *Job) fuzzCycle(ctx context.Context, cycle int) error {
	sc := SubjectCycle{Name: j.Subject.Name, Cycle: cycle}
	spaths := j.Pathset.SubjectPaths(sc)
	if err := j.Driver.FuzzSingle(ctx, j.Rng.Int31(), j.Subject.Litmus, spaths.Litmus, spaths.Trace); err != nil {
		return err
	}

	nsub := j.fuzzedSubject(sc, spaths)
	return j.sendSubject(ctx, nsub)
}

// fuzzedSubject makes a copy of this Job's subject with the cycled name sc and fuzz fileset spaths.
func (j *Job) fuzzedSubject(sc SubjectCycle, spaths subject.FuzzFileset) subject.Named {
	nsub := j.Subject
	nsub.Name = sc.String()
	nsub.Fuzz = &spaths
	return nsub
}

// sendSubject tries to send subject s down this Job's results channel.
func (j *Job) sendSubject(ctx context.Context, s subject.Named) error {
	select {
	case j.ResCh <- j.builderReq(s):
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// builderReq makes a builder request for adding s to the fuzzed corpus.
func (j *Job) builderReq(s subject.Named) corpus.BuilderReq {
	return corpus.BuilderReq{
		Name: s.Name,
		Req:  corpus.AddReq(s.Subject),
	}
}
