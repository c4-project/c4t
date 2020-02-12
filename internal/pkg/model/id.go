package model

import (
	"errors"
	"fmt"
	"strings"
)

const idSep = '.'

var (
	// ErrNoTags occurs when one calls NewId with no tags.
	ErrNoTags = errors.New("no tags")

	// ErrTagHasSep occurs when a tag passed to NewId contains the separator rune.
	ErrTagHasSep = errors.New("tag contains separator")

	// EmptyId is the empty ACT ID.
	EmptyId = Id{}
)

// Id represents an ACT ID.
type Id struct {
	tags []string
}

// Tags extracts the tags comprising an ID as a slice.
func (i Id) Tags() []string {
	return i.tags
}

// String converts an ACT ID to a string.
func (i Id) String() string {
	return strings.Join(i.tags, string(idSep))
}

// NewId tries to construct an ACT ID from tags.
// It fails if any of the tags contain a separator.
// It also fails if no tags are passed.
func NewId(tags ...string) (Id, error) {
	// We could use the signature (tag string, tags ...string) to enforce the 'at least one segment' rule,
	// but this'd make it harder to splat in a []string.
	if err := validateTags(tags); err != nil {
		return Id{nil}, fmt.Errorf("tag validation failed for %v: %w", tags, err)
	}

	// Normalise the empty tag.
	if len(tags) == 1 && tags[0] == "" {
		return EmptyId, nil
	}

	return Id{tags}, nil
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

// IdFromString converts a string to an ACT ID.
func IdFromString(s string) Id {
	return Id{strings.Split(s, ".")}
}

// MarshalText implements text marshalling for IDs by stringifying them.
func (i Id) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements text unmarshalling for IDs by unstringifying them.
func (i *Id) UnmarshalText(b []byte) error {
	*i = IdFromString(string(b))
	return nil
}
