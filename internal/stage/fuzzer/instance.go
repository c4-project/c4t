// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"context"
	"errors"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/subject"
)

// Instance contains state for a single fuzzer instance.
type Instance struct {
	// Normalise contains the subject for which this Instance is responsible.
	Subject subject.Named

	// Driver is the low-level fuzzer driver set.
	Driver Driver

	// SubjectCycles is the number of times each subject should be fuzzed.
	SubjectCycles int

	// Machine is the machine, if any, being targeted by the fuzzer.
	// Knowledge of the machine can be used to shape things like thread counts.
	Machine *machine.Machine

	// Config is the specific configuration, if any, for the fuzzer.
	Config *fuzzer.Configuration

	// Pathset points to the pathset to use to work out where to store fuzz output.
	Pathset SubjectPather

	// Rng is the random number generator to use for fuzz seeds.
	Rng *rand.Rand

	// ResCh is the channel onto which each fuzzed subject should be sent.
	ResCh chan<- builder.Request
}

// Fuzz performs a single fuzzing instance.
func (j *Instance) Fuzz(ctx context.Context) error {
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
func (j *Instance) check() error {
	if j.Driver == nil {
		return ErrDriverNil
	}
	if j.Pathset == nil {
		return iohelp.ErrPathsetNil
	}
	if j.Rng == nil {
		return errors.New("RNG nil")
	}
	return nil
}

func (j *Instance) fuzzCycle(ctx context.Context, cycle int) error {
	sc := SubjectCycle{Name: j.Subject.Name, Cycle: cycle}
	jb := j.makeJob(sc)

	stime := time.Now()
	if err := j.Driver.Fuzz(ctx, jb); err != nil {
		return err
	}
	dur := time.Since(stime)

	// TODO(@MattWindsor91): should we double-check the architecture here?
	l, err := litmus.New(jb.OutLitmus, litmus.WithArch(id.ArchC), litmus.PopulateStatsFrom(ctx, j.Driver))
	if err != nil {
		return nil
	}

	fz := subject.Fuzz{
		Duration: dur,
		Litmus:   *l,
		Trace:    jb.OutTrace,
	}

	nsub := j.fuzzedSubject(sc, &fz)
	return builder.AddRequest(&nsub).SendTo(ctx, j.ResCh)
}

func (j *Instance) makeJob(sc SubjectCycle) fuzzer.Job {
	jb := fuzzer.Job{
		Seed:    j.Rng.Int31(),
		In:      j.Subject.Source.Path,
		Machine: j.Machine,
		Config:  j.Config,
		// These two are filepaths, but the job stores slashpaths.
		OutLitmus: filepath.ToSlash(j.Pathset.SubjectLitmus(sc)),
		OutTrace:  filepath.ToSlash(j.Pathset.SubjectTrace(sc)),
	}
	return jb
}

// fuzzedSubject makes a copy of this Instance's subject with the cycled name sc and fuzz fz.
func (j *Instance) fuzzedSubject(sc SubjectCycle, fz *subject.Fuzz) subject.Named {
	nsub := j.Subject
	nsub.Name = sc.String()
	nsub.Fuzz = fz
	return nsub
}
