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

// ReadHeader tries to read a Header from JSON in rd.
func ReadHeader(rd io.Reader) (*Header, error) {
	hdr := Header{}
	dec := json.NewDecoder(rd)
	err := dec.Decode(&hdr)
	return &hdr, err
}
