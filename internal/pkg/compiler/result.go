package compiler

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// Result logs information about an attempt to compile a subject with a compiler under test.
type Result struct {
	// BackendID is the CompilerID of the backend that generated the harness compiled by this compiler.
	BackendID model.ID

	// CompilerID is the CompilerID of the compiler used to compile this binary.
	CompilerID model.ID

	// SubName is the name of the test subject.
	SubName string

	// Success gets whether the compilation succeeded (possibly with errors).
	Success bool

	// PathBin, on success, provides the path to the compiled binary.
	PathBin string

	// PathLog provides the path to the compiler's stderr log.
	PathLog string
}
