// Package act_tester_plan contains the app-specific logic for the test planner.
package act_tester_plan

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

var (
	ErrNoCorpus = errors.New("no test files supplied")
)

// Planner runs the test planner p.
func (p *Planner) Plan() error {
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

	return dumpPlan(&plan)
}

// dumpPlan dumps plan p to stdout.
func dumpPlan(p *model.Plan) error {
	// TODO(@MattWindsor91): output to other files
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	return enc.Encode(p)
}
