package model

import (
	"errors"
	"fmt"
	"strings"
)

const idSep = '.'

var (
	// ErrNoTags occurs when one calls NewID with no tags.
	ErrNoTags = errors.New("no tags")

	// ErrTagHasSep occurs when a tag passed to NewID contains the separator rune.
	ErrTagHasSep = errors.New("tag contains separator")

	// EmptyID is the empty ACT ID.
	EmptyID = ID{}
)

// ID represents an ACT ID.
type ID struct {
	tags []string
}

// Tags extracts the tags comprising an ID as a slice.
func (i ID) Tags() []string {
	return i.tags
}

// String converts an ACT ID to a string.
func (i ID) String() string {
	return strings.Join(i.tags, string(idSep))
}

// NewID tries to construct an ACT ID from tags.
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
		return EmptyID, nil
	}

	return ID{tags}, nil
}

func validateTags(tags []string) error {
	if len(tags) == 0 {
		return ErrNoTags
	}
	for _, t := range tags {
		if strings.ContainsRune(t, idSep) {
			return ErrTagHasSep
		}
	}
	return nil
}

// IDFromString converts a string to an ACT ID.
func IDFromString(s string) ID {
	return ID{strings.Split(s, ".")}
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
