// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp

import "fmt"

// PluralQuantity formats a quantity n using the components stem, one, and many as follows:
// if n is 1, we return '[n] [stem][one]'; else, '[n] [stem][many]'.
func PluralQuantity(n int, stem, one, many string) string {
	return fmt.Sprintf("%d %s%s", n, stem, pluralSuffix(n == 1, one, many))
}

// pluralSuffix selects one if isOne is true, else many.
func pluralSuffix(isOne bool, one, many string) string {
	if isOne {
		return one
	}
	return many
}
