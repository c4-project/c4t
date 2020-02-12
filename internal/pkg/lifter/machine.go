package lifter

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// machine contains state used in a single-machine lifting job.
type machine struct {
	// Corpus is the corpus we are lifting.
	Corpus model.Corpus

	// Dir is the directory to which lifted harnesses should be created.
	Dir string

	// Machine is the plan for this machine.
	Machine model.MachinePlan

	// Maker is the parent harness maker.
	Maker HarnessMaker

	// ResCh is the channel to which the machine lifter should send lifting results.
	ResCh chan<- result
}

func (m *machine) lift(ctx context.Context) error {
	logrus.WithField("machine", m.Machine.Id).Debugln("lifting machine")

	for _, a := range m.Machine.Arches() {
		if err := m.liftArch(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

func (m *machine) liftArch(ctx context.Context, arch model.Id) error {
	dir, derr := buildAndMkDir(m.Dir, arch.Tags()...)
	if derr != nil {
		return derr
	}

	for i := range m.Corpus {
		if err := m.liftSubject(ctx, arch, dir, &(m.Corpus[i])); err != nil {
			return err
		}
	}
	return nil
}

func (m *machine) liftSubject(ctx context.Context, arch model.Id, dir string, s *model.Subject) error {
	dir, derr := buildAndMkDir(dir, s.Name)
	if derr != nil {
		return derr
	}

	spec := model.HarnessSpec{
		Backend: m.Machine.Backend.Id,
		Arch:    arch,
		InFile:  s.Litmus,
		OutDir:  dir,
	}

	logrus.WithField("spec", spec).Debugln("making harness")
	files, err := m.Maker.MakeHarness(spec)
	if err != nil {
		return err
	}

	res := result{
		Arch:    arch,
		Harness: model.Harness{Dir: dir, Files: files},
		Machine: m.Machine.Id,
		Subject: s,
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
