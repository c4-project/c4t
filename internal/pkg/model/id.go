package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const idSep = '.'

var (
	// ErrTagHasSep occurs when one calls NewId with no tags.
	ErrNoTags = errors.New("no tags")

	// ErrTagHasSep occurs when a tag passed to NewId contains the separator rune.
	ErrTagHasSep = errors.New("tag contains separator")
)

// Id represents an ACT ID.
type Id struct {
	tags []string
}

// String converts an ACT ID to a string.
func (i Id) String() string {
	return strings.Join(i.tags, string(idSep))
}

// NewId tries to construct an ACT ID from tags.
// It fails if any of the tags contain a separator.
// It also fails if no tags are passed.
func IdFromTag(tags ...string) (Id, error) {
	// We could use the signature (tag string, tags ...string) to enforce the 'at least one segment' rule,
	// but this'd make it harder to splat in a []string.
	if err := validateTags(tags); err != nil {
		return Id{nil}, fmt.Errorf("tag validation failed for %v: %w", tags, err)
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

// MarshalJSON implements JSON marshalling for IDs by stringifying them.
func (i Id) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *Id) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s != nil {
		*i = IdFromString(*s)
	}
	return nil
}
