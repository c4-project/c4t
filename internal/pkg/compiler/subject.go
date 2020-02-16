package compiler

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// SubjectCompile describes the unique name of a particular instance of the batch compiler.
type SubjectCompile struct {
	// Name is the name of the subject.
	Name string

	// CompilerID is the ID of the compiler.
	CompilerID model.ID
}
