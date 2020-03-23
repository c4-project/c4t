// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import "io"

// CopyClose copies src to dst, then closes both files.
func CopyClose(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	n, cerr := io.Copy(dst, src)
	derr := dst.Close()
	serr := src.Close()

	return n, FirstError(cerr, derr, serr)
}
