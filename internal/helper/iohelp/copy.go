// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import "io"

// CopyCloseSrc copies src to dst, then closes src.
func CopyCloseSrc(dst io.Writer, src io.ReadCloser) (int64, error) {
	n, cerr := io.Copy(dst, src)
	serr := src.Close()
	return n, FirstError(cerr, serr)
}

// CopyClose copies src to dst, then closes both files.
func CopyClose(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	n, cerr := CopyCloseSrc(dst, src)
	derr := dst.Close()

	return n, FirstError(cerr, derr)
}
