// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
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
	ResCh chan<- builder.Request
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

	stime := time.Now()
	if err := j.Driver.FuzzSingle(ctx, j.Rng.Int31(), j.Subject.Litmus, spaths); err != nil {
		return err
	}
	fz := subject.Fuzz{
		Duration: time.Since(stime),
		Files:    spaths,
	}

	nsub := j.fuzzedSubject(sc, &fz)
	return builder.AddRequest(&nsub).SendTo(ctx, j.ResCh)
}

// fuzzedSubject makes a copy of this Job's subject with the cycled name sc and fuzz fileset spaths.
func (j *Job) fuzzedSubject(sc SubjectCycle, fz *subject.Fuzz) subject.Named {
	nsub := j.Subject
	nsub.Name = sc.String()
	nsub.Fuzz = fz
	return nsub
}
