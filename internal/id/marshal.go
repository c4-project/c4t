// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

// MarshalText implements text marshalling for IDs by stringifying them.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements text unmarshalling for IDs by unstringifying them.
func (i *ID) UnmarshalText(b []byte) error {
	return i.Set(string(b))
}
