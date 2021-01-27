// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"errors"
	"fmt"
	"strings"
)

// Tag classifies a state line.
type Tag int

// BadTag occurs when we try to unmarshal a string that doesn't correspond to a Tag.
var BadTag = errors.New("bad tag name")

const (
	// TagUnknown represents a state that is not known to be either a witness or a counter-example.
	TagUnknown Tag = iota // unknown
	// TagWitness represents a state that validates a condition.
	TagWitness // witness
	// TagCounter represents a state that goes against a condition.
	TagCounter // counter
	// TagLast refers to the last tag.
	TagLast = TagCounter
)

//go:generate stringer -type Tag -linecomment

// MarshalText marshals a Tag into text.
func (i Tag) MarshalText() (text []byte, err error) {
	return []byte(i.String()), nil
}

// UnmarshalText unmarshals text into a Tag.
func (i *Tag) UnmarshalText(text []byte) error {
	tstr := string(text)
	for *i = TagUnknown; *i <= TagLast; *i++ {
		if strings.EqualFold(tstr, i.String()) {
			return nil
		}
	}
	return fmt.Errorf("%w: %q", BadTag, tstr)
}
