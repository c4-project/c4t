// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package status

import "encoding/json"

// MarshalText marshals a Status to text via its string representation.
func (s Status) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText unmarshals a Status from text via its string representation.
func (s *Status) UnmarshalText(text []byte) error {
	var err error
	*s, err = OfString(string(text))
	return err
}

// MarshalJSON marshals a Status to JSON via its string representation.
func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnarshalJSON unmarshals a Status from JSON via its string representation.
func (s *Status) UnmarshalJSON(bytes []byte) error {
	var (
		sstr string
		err  error
	)
	if err = json.Unmarshal(bytes, &sstr); err != nil {
		return err
	}
	*s, err = OfString(sstr)
	return err
}
