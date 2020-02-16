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

	// SepQual is a separator used to distinguish two parts of a qualified ID when there should be no ambiguity.
	// It is exported for testing and sanitisation purposes.
	SepQual = ':'

	// seps contains tagSep and qualSep.
	seps = ".:"
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
		if strings.ContainsAny(t, seps) {
			return fmt.Errorf("%w: tag %q", ErrTagHasSep, t)
		}
	}
	return nil
}

// IDFromString converts a string to an ACT ID.
func IDFromString(s string) ID {
	return ID{strings.Split(s, string(SepTag))}
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

// MachQualID is a type for IDs of things qualified by machine IDs.
type MachQualID struct {
	// MachineID is the ID of the qualifying machine.
	MachineID ID

	// ID is the ID of the qualified item.
	ID ID
}

// String converts a machine-qualified ID to a string.
//
// It does so by joining the two IDs together with a distinct separator from the usual tag separator;
// this is to prevent ambiguity.
//
// To combine the two IDs into one 'fully qualified ID', use FQID instead.
func (m MachQualID) String() string {
	strs := []string{m.MachineID.String(), m.ID.String()}
	return strings.Join(strs, string(SepQual))
}

// FQID converts this machine-qualified ID into a single fully-qualified ID.
func (m MachQualID) FQID() ID {
	return m.MachineID.Join(m.ID)
}
