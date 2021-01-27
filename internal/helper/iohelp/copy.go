// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"io"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

// CopyCloseSrc copies src to dst, then closes src.
// It closes src even if there was an error while copying.
func CopyCloseSrc(dst io.Writer, src io.ReadCloser) (int64, error) {
	n, cerr := io.Copy(dst, src)
	serr := src.Close()
	return n, errhelp.FirstError(cerr, serr)
}

// CopyClose copies src to dst, then closes both.
// It closes src and dst even if there was an error while copying.
func CopyClose(dst io.WriteCloser, src io.ReadCloser) (int64, error) {
	n, cerr := CopyCloseSrc(dst, src)
	derr := dst.Close()

	return n, errhelp.FirstError(cerr, derr)
}
