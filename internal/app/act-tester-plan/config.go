package act_tester_plan

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
)

// Config represents the input config to the test planner.
type Config struct {
	// Filter is the compiler filter to use to select compilers to test.
	Filter interop.CompilerFilter

	// Corpus is a list of paths to files that form the incoming test corpus.
	Corpus []string
}
