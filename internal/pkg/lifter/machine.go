// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"

	"github.com/sirupsen/logrus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// machine contains state used in a single-machine lifting job.
type machine struct {
	// Corpus is the corpus we are lifting.
	Corpus corpus.Corpus

	// Dir is the directory to which lifted harnesses should be created.
	Dir string

	// MachineID is the CompilerID for this machine.
	MachineID model.ID

	// Machine is the plan for this machine.
	Machine plan.MachinePlan

	// Maker is the parent harness maker.
	Maker HarnessMaker

	// ResCh is the channel to which the machine lifter should send lifting results.
	ResCh chan<- result
}

func (m *machine) lift(ctx context.Context) error {
	logrus.WithField("machine", m.MachineID).Debugln("lifting machine")
	for _, a := range m.Machine.Arches() {
		if err := m.liftArch(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

func (m *machine) liftArch(ctx context.Context, arch model.ID) error {
	dir, derr := buildAndMkDir(m.Dir, arch.Tags()...)
	if derr != nil {
		return derr
	}

	return m.Corpus.Each(func(s subject.Named) error {
		return m.liftSubject(ctx, arch, dir, &s)
	})
}

func (m *machine) liftSubject(ctx context.Context, arch model.ID, dir string, s *subject.Named) error {
	dir, derr := buildAndMkDir(dir, s.Name)
	if derr != nil {
		return derr
	}

	path, perr := s.BestLitmus()
	if perr != nil {
		return perr
	}

	spec := model.HarnessSpec{
		Backend: m.Machine.Backend.ID,
		Arch:    arch,
		InFile:  path,
		OutDir:  dir,
	}

	logrus.WithField("spec", spec).Debugln("making harness")
	files, err := m.Maker.MakeHarness(ctx, spec)
	if err != nil {
		return err
	}

	res := result{
		MArch:   model.MachQualID{MachineID: m.MachineID, ID: arch},
		Harness: subject.Harness{Dir: dir, Files: files},
		Subject: s.Name,
	}

	return m.sendResult(ctx, res)
}

func (m *machine) sendResult(ctx context.Context, r result) error {
	select {
	case m.ResCh <- r:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
