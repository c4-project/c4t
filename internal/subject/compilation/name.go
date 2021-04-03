// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"fmt"
	"path"

	"github.com/c4-project/c4t/internal/id"
)

// Name describes the unique name of a particular instance of the batch compiler.
type Name struct {
	// SubjectName is the name of the subject.
	SubjectName string

	// CompilerID is the ID of the compiler.
	CompilerID id.ID
}

// String gets a stringified version of the name.
func (n Name) String() string {
	return fmt.Sprintf("%s@%s", n.SubjectName, n.CompilerID)
}

// Path gets a slashpath fragment that can be used to locate this compilation unambiguously in a directory tree.
func (n Name) Path() string {
	return path.Join(append(n.CompilerID.Tags(), n.SubjectName)...)
}
