// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// SubjectCompile describes the unique name of a particular instance of the batch compiler.
type SubjectCompile struct {
	// Name is the name of the subject.
	Name string

	// CompilerID is the ID of the compiler.
	CompilerID model.ID
}
