// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"errors"
	"fmt"
)

// TagGlob is the tag used in a glob expression to indicate that everything before it should be a prefix,
// and everything after a suffix, of the matched ID.
const TagGlob = "*"

// ErrBadGlob occurs when Matches gets a malformed glob expression.
var ErrBadGlob = errors.New("malformed glob expression")

// HasPrefix tests whether prefix is a prefix of this ID.
func (i ID) HasPrefix(prefix ID) bool {
	li := len(i.tags)
	lp := len(prefix.tags)
	switch {
	case li < lp:
		return false
	case li == lp:
		return i.Equal(prefix)
	default:
		return i.slice(0, len(prefix.tags)).Equal(prefix)
	}
}

// HasPrefix tests whether suffix is a suffix of this ID.
func (i ID) HasSuffix(suffix ID) bool {
	li := len(i.tags)
	ls := len(suffix.tags)
	switch {
	case li < ls:
		return false
	case li == ls:
		return i.Equal(suffix)
	default:
		return i.slice(len(i.tags)-len(suffix.tags), len(i.tags)).Equal(suffix)
	}
}

// Matches tests whether this ID matches the glob ID expression glob.
// glob should be either a literal ID, or an ID with exactly one tag equal to TagGlob.
func (i ID) Matches(glob ID) (bool, error) {
	prefix, suffix, exact, err := split(glob)
	if err != nil {
		return false, err
	}
	if exact {
		return i.Equal(glob), nil
	}
	return i.HasPrefix(prefix) && i.HasSuffix(suffix), nil
}

func (i ID) slice(from int, to int) ID {
	return ID{tags: i.tags[from:to]}
}

// split splits a glob ID into a prefix, suffix, and error.
// If there is no glob character, we return a flag that specifies that the match should be exact.
func split(glob ID) (prefix ID, suffix ID, exact bool, err error) {
	globIndex := -1
	for i, tag := range glob.tags {
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

	prefix.tags = glob.tags[:globIndex]
	suffix.tags = glob.tags[globIndex+1:]
	return prefix, suffix, false, err
}
