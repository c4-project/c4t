// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import "encoding/json"

// MarshalText implements text marshalling for IDs by stringifying them.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements text unmarshalling for IDs by unstringifying them.
func (i *ID) UnmarshalText(b []byte) error {
	var err error
	*i, err = TryFromString(string(b))
	return err
}

// MarshalJSON marshals an ID to JSON via its string representation.
func (i ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnarshalJSON unmarshals an ID from JSON via its string representation.
func (i *ID) UnmarshalJSON(bytes []byte) error {
	var (
		sstr string
		err  error
	)
	if err = json.Unmarshal(bytes, &sstr); err != nil {
		return err
	}
	*i, err = TryFromString(sstr)
	return err
}
