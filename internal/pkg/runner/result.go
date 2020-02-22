package runner

import "time"

// Result is the type of results from a single test run.
type Result struct {
	// Start is the time at which the run commenced.
	Start time.Time

	// Compilers is the set of results that were reported in this test run.
	// Each key is a stringified form of a compiler CompilerID.
	Compilers map[string]CompilerResult
}

// CompilerResult contains results from a single compiler.
type CompilerResult struct {
}

type SubjectResult struct {
	// RawObs is the path to this subject's raw observation file.
	RawObs string `toml:"log,omitempty"`
	// Obs is the path to this subject's processed observation file.
	Obs string `toml:"log,omitempty"`
}
