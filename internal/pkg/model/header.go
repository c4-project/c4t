package model

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
