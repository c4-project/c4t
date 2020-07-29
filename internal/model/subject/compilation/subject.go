// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Name describes the unique name of a particular instance of the batch compiler.
type Name struct {
	// SubjectName is the name of the subject.
	SubjectName string

	// CompilerID is the ID of the compiler.
	CompilerID id.ID
}
