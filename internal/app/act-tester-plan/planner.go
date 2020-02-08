// Package act_tester_plan contains the app-specific logic for the test planner.
package act_tester_plan

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	// Source is the planner's information source.
	Source interface {
		interop.BackendFinder
		interop.CompilerLister
	}

	// Filter is the compiler filter to use to select compilers to test.
	Filter model.CompilerFilter

	// CorpusSize is the requested size of the test corpus.
	// If zero, no corpus sampling is done, but the planner will still error if the final corpus size is 0.
	// If nonzero, the corpus will be sampled if larger than the size, and an error occurs if the final size is below
	// that requested.
	CorpusSize int

	// Corpus is a list of paths to files that form the incoming test corpus.
	Corpus []string
}

// Planner runs the test planner p.
func (p *Planner) Plan() error {
	// Early out to prevent us from doing any planning if we received no files.
	if len(p.Corpus) == 0 {
		return ErrNoCorpus
	}

	var plan model.Plan

	rand.Seed(time.Now().UnixNano())
	plan.Init()

	var err error
	if plan.Machines, err = p.planMachines(); err != nil {
		return err
	}

	if plan.Corpus, err = p.planCorpus(plan.Seed); err != nil {
		return err
	}

	return dumpPlan(&plan)
}

// dumpPlan dumps plan p to stdout.
func dumpPlan(p *model.Plan) error {
	// TODO(@MattWindsor91): output to other files
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
