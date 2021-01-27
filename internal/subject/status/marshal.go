// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package status

import "encoding/json"

// MarshalText marshals a Status to text via its string representation.
func (i Status) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText unmarshals a Status from text via its string representation.
func (i *Status) UnmarshalText(text []byte) error {
	var err error
	*i, err = FromString(string(text))
	return err
}

// MarshalJSON marshals a Status to JSON via its string representation.
func (i Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnarshalJSON unmarshals a Status from JSON via its string representation.
func (i *Status) UnmarshalJSON(bytes []byte) error {
	var (
		sstr string
		err  error
	)
	if err = json.Unmarshal(bytes, &sstr); err != nil {
		return err
	}
	*i, err = FromString(sstr)
	return err
}
