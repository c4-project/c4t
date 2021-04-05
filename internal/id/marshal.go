// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

import (
	"errors"

	"github.com/pelletier/go-toml"
)

// MarshalText implements text marshalling for IDs by stringifying them.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements text unmarshalling for IDs by unstringifying them.
func (i *ID) UnmarshalText(b []byte) error {
	return i.Set(string(b))
}

func (i ID) MarshalTOML() ([]byte, error) {
	return toml.Marshal(i.String())
}

func (i *ID) UnmarshalTOML(b interface{}) error {
	str, ok := b.(string)
	if !ok {
		return errors.New("IDs can only be unmarshalled from strings")
	}
	return i.Set(str)
}
