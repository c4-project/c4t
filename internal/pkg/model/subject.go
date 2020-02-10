package model

// Subject represents a single test subject in a corpus.
type Subject struct {
	// Name is the name of this
	Name string

	// Litmus is the path to this subject's current Litmus file.
	Litmus string

	// OrigLitmus is the path to this subject's original Litmus file.
	// If empty, then Litmus is the original file.
	OrigLitmus string

	// TracePath is the path to this subject's fuzzer trace file.
	// If empty, this subject hasn't been fuzzed by act-tester-fuzz.
	TracePath string

	// HarnessPaths contains the paths of every file in this subject's test harness.
	// If nil, this subject hasn't had a harness generated.
	HarnessPaths map[string][]string
}
