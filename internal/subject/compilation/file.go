// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"errors"

	"github.com/MattWindsor91/c4t/internal/subject/normpath"

	"github.com/1set/gut/ystring"
)

var (
	// ErrNoCompilerLog occurs when we ask for the compiler log of a subject that doesn't have one.
	ErrNoCompilerLog = errors.New("compiler result has no log file")
)

// ReadLog tries to read in the log for compiler, taking paths relative to root.
// If the compiler log doesn't exist relative to root, and its path is of the form FOO/BAR, we assume that it is in a
// saved tarball called FOO.tar.gz (as file BAR) in root, and attempt to extract it.
func (c *CompileFileset) ReadLog(root string) ([]byte, error) {
	if ystring.IsBlank(c.Log) {
		return nil, ErrNoCompilerLog
	}
	return normpath.ReadSubjectFile(root, c.Log)
}
