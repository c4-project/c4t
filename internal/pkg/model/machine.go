package model

// Machine represents the information about a machine that is relevant to the tester.
type Machine struct {
	// ID is the identifier of the machine.
	ID ID `toml:"id"`

	// Cores is the number of known cores on the machine.
	// If zero, there is no known core count.
	Cores int `toml:"cores,omitzero"`
}
