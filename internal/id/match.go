// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"errors"
	"fmt"
	"strings"
)

// TagGlob is the tag used in a glob expression to indicate that everything before it should be a prefix,
// and everything after a suffix, of the matched ID.
const TagGlob = "*"

// ErrBadGlob occurs when Matches gets a malformed glob expression.
var ErrBadGlob = errors.New("malformed glob expression")

// HasPrefix tests whether prefix is a prefix of this ID.
func (i ID) HasPrefix(prefix ID) bool {
	switch {
	// Easiest cases can be done by direct string analysis.
	case prefix.IsEmpty():
		return true
	case len(i.repr) < len(prefix.repr):
		return false
	case len(i.repr) == len(prefix.repr):
		return i.Equal(prefix)
	default:
		// We can use string prefixing, with a caveat: we must ensure the last tag of prefix is contained verbatim in i.
		// Because i is larger, we assume it has more tags, so we can do this by adding SepTag.
		return strings.HasPrefix(i.repr, prefix.repr+SepTag)
	}
}

// HasPrefix tests whether suffix is a suffix of this ID.
func (i ID) HasSuffix(suffix ID) bool {
	// As HasPrefix, but with the final case flipped around for suffixing.
	switch {
	case suffix.IsEmpty():
		return true
	case len(i.repr) < len(suffix.repr):
		return false
	case len(i.repr) == len(suffix.repr):
		return i.Equal(suffix)
	default:
		return strings.HasSuffix(i.repr, SepTag+suffix.repr)
	}
}

// Matches tests whether this ID matches the glob ID expression glob.
// glob should be either a literal ID, or an ID with exactly one tag equal to TagGlob.
func (i ID) Matches(glob ID) (bool, error) {
	prefix, suffix, exact, err := splitGlob(glob)
	if err != nil {
		return false, err
	}
	if exact {
		return i.Equal(glob), nil
	}
	return i.HasPrefix(prefix) && i.HasSuffix(suffix), nil
}

// splitGlob splits a glob ID into a prefix, suffix, and error.
// If there is no glob character, we return a flag that specifies that the match should be exact.
func splitGlob(glob ID) (prefix ID, suffix ID, exact bool, err error) {
	// TODO(@MattWindsor91): use a string based algorithm here?

	globIndex := -1
	gtags := glob.Tags()
	for i, tag := range gtags {
		if tag == TagGlob {
			if globIndex != -1 {
				return prefix, suffix, exact, fmt.Errorf("%w: more than one '*' character", ErrBadGlob)
			}
			globIndex = i
		}
	}
	if globIndex == -1 {
		return prefix, suffix, true, nil
	}

	prefix = unsafeJoin(gtags[:globIndex]...)
	suffix = unsafeJoin(gtags[globIndex+1:]...)
	return prefix, suffix, false, nil
}
