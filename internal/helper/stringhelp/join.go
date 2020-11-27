// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp

import (
	"strings"

	"github.com/1set/gut/ystring"
)

// JoinNonEmpty joins the non-empty strings in xs using sep.
func JoinNonEmpty(sep string, xs ...string) string {
	nxs := make([]string, 0, len(xs))
	for _, x := range xs {
		if ystring.IsNotEmpty(x) {
			nxs = append(nxs, x)
		}
	}
	return strings.Join(nxs, sep)
}
