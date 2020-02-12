// Package planner contains the logic for the test planner.
package planner

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BackendFinder is the interface of things that can find backends for machines.
type BackendFinder interface {
	// FindBackend asks for a backend with the given style on any one of machines,
	// or a default machine if none have such a backend.
	FindBackend(style model.ID, machines ...model.ID) (*model.Backend, error)
}

// Source is the composite interface of types that can provide the requisite information a Planner needs about
// backends, compilers, and subjects.
type Source interface {
	BackendFinder
	CompilerLister
	SubjectProber
}

// Planner holds all configuration for the test planner.
type Planner struct {
	// Source is the planner's information source.
	Source Source

	// Filter is the compiler filter to use to select compilers to test.
	Filter model.CompilerFilter

	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int

	// InFiles is a list of paths to files that form the incoming test corpus.
	InFiles []string
}

// Plan runs the test planner p.
func (p *Planner) Plan() error {
	// Early out to prevent us from doing any planning if we received no files.
	if len(p.InFiles) == 0 {
		return model.ErrNoCorpus
	}

	var plan model.Plan

	rand.Seed(time.Now().UnixNano())
	plan.Init()

	logrus.Infoln("Planning machines...")
	var err error
	if plan.Machines, err = p.planMachines(); err != nil {
		return err
	}

	logrus.Infoln("Planning corpus...")
	if plan.Corpus, err = p.planCorpus(plan.Seed); err != nil {
		return err
	}

	return plan.Dump()
}
