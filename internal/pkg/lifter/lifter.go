// Package lifter contains the part of the tester framework that lifts litmus tests to compilable C.
// It does so by means of a backend HarnessMaker.
package lifter

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// HarnessMaker is an interface capturing the ability to make test harnesses.
type HarnessMaker interface {
	// MakeHarness asks the harness maker to make the test harness described by spec.
	// It returns a list outfiles of files created (C files, header files, etc.), and/or an error err.
	MakeHarness(spec model.HarnessSpec) (outFiles []string, err error)
}

// Lifter holds the main configuration for the lifter part of the tester framework.
type Lifter struct {
	// Maker is a harness maker.
	Maker HarnessMaker
}

func (l *Lifter) Lift() error {
	// TODO(@MattWindsor91)
	return nil
}
