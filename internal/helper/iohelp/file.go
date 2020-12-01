// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"os"
)

// IsFileEmpty checks whether the open file f is empty.
// It can fail if we can't stat the file.
func IsFileEmpty(f *os.File) (bool, error) {
	fi, err := f.Stat()
	if err != nil {
		return false, err
	}
	return fi.Size() == 0, nil
}
