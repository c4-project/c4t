package model

// HarnessSpec is a specification of how to make a test harness.
type HarnessSpec struct {
	// Backend is the fully-qualified identifier of the backend to use to make this harness.
	Backend Id

	// Emits is the ID of the architecture for which a harness should be prepared.
	Emits Id

	// InFile is the path to the input Litmus test file.
	InFile string

	// OutDir is the path to the output Litmus
	OutDir string
}
