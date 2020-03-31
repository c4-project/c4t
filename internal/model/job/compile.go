// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

import (
	"github.com/MattWindsor91/act-tester/internal/model/compiler"
)

// Compile represents a request to compile a list of files to an executable given a particular compiler.
type Compile struct {
	// Compiler describes the compiler to use for the compilation.
	Compiler *compiler.Compiler

	// In is the list of files to be sent to the compiler.
	In []string
	// Out is the file to be received from the compiler.
	Out string
}
