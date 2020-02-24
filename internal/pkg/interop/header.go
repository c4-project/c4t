// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package interop

import (
	"encoding/json"
	"io"
)

// Header represents a Litmus test header in the form that act-c accepts and dumps.
type Header struct {
	// Name is the name of the Litmus test.
	Name string `json:"name"`

	// Locations is the list of locations present in the Litmus test.
	Locations []string `json:"locations"`

	// Init is the initialiser block for the Litmus test.
	Init map[string]int `json:"init"`

	// Postcondition is the Litmus postcondition.
	Postcondition string `json:"postcondition"`
}

// ReadHeader tries to read a Header from JSON in r.
func (h *Header) Read(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&h)
}
