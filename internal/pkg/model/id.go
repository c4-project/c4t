// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// SepTag is the identifier tag separator.
	// It is exported for testing and sanitisation purposes.
	SepTag = '.'
)

var (
	// ErrNoTags occurs when one calls NewID with no tags.
	ErrNoTags = errors.New("no tags")

	// ErrTagHasSep occurs when a tag passed to NewID contains the separator rune.
	ErrTagHasSep = errors.New("tag contains separator")
)

// ID represents an ACT ID.
type ID struct {
	tags []string
}

// NewID tries to construct an ACT CompilerID from tags.
// It fails if any of the tags contain a separator.
// It also fails if no tags are passed.
func NewID(tags ...string) (ID, error) {
	// We could use the signature (tag string, tags ...string) to enforce the 'at least one segment' rule,
	// but this'd make it harder to splat in a []string.
	if err := validateTags(tags); err != nil {
		return ID{nil}, fmt.Errorf("tag validation failed for %v: %w", tags, err)
	}

	// Normalise the empty tag.
	if len(tags) == 1 && tags[0] == "" {
		return ID{}, nil
	}

	return ID{tags}, nil
}

func validateTags(tags []string) error {
	if len(tags) == 0 {
		return ErrNoTags
	}
	for _, t := range tags {
		if strings.ContainsRune(t, SepTag) {
			return fmt.Errorf("%w: tag %q", ErrTagHasSep, t)
		}
	}
	return nil
}

// TryIDFromString tries to convert a string to an ACT ID.
// It returns any validation error arising.
func TryIDFromString(s string) (ID, error) {
	return NewID(strings.Split(s, string(SepTag))...)
}

// IDFromString converts a string to an ACT ID.
// It returns the empty ID if there is an error.
func IDFromString(s string) ID {
	id, err := TryIDFromString(s)
	if err != nil {
		return ID{}
	}
	return id
}

// IsEmpty gets whether this ID is empty.
func (i ID) IsEmpty() bool {
	return len(i.tags) == 0
}

// Tags extracts the tags comprising an ID as a slice.
func (i ID) Tags() []string {
	return i.tags
}

// String converts an ACT ID to a string.
func (i ID) String() string {
	return strings.Join(i.tags, string(SepTag))
}

// Join appends r to this ID, creating a new ID.
func (i ID) Join(r ID) ID {
	if i.IsEmpty() {
		return r
	}
	if r.IsEmpty() {
		return i
	}
	return ID{append(i.tags, r.tags...)}
}

// MarshalText implements text marshalling for IDs by stringifying them.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements text unmarshalling for IDs by unstringifying them.
func (i *ID) UnmarshalText(b []byte) error {
	*i = IDFromString(string(b))
	return nil
}

// Less compares two IDs lexicographically.
func (i ID) Less(i2 ID) bool {
	for j := 0; j < len(i.tags) && j < len(i2.tags); j++ {
		switch {
		case i.tags[j] < i2.tags[j]:
			return true
		case i.tags[j] > i2.tags[j]:
			return false
		}
	}
	return len(i.tags) < len(i2.tags)
}
