package act_tester_plan

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Planner holds all configuration for the test planner.
type Planner struct {
	// Act tells the planner how to run ACT.
	Act interop.ActRunner

	// Filter is the compiler filter to use to select compilers to test.
	Filter model.CompilerFilter

	// Corpus is a list of paths to files that form the incoming test corpus.
	Corpus []string
}
