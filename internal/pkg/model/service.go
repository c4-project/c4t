package model

// Service is a structure collecting information about a 'service' (a compiler or backend).
type Service struct {
	// Id is the (likely unqualified) ACT ID of the service.
	Id Id `json:"id"`

	// IdQualified is true if Id is qualified by the machine ID.
	IdQualified bool `json:"id_qualified,omitempty"`

	// MachineId is the ACT ID of the service's parent machine.
	// It may be empty if there is no need to track it.
	MachineId *Id `json:"machine_id,omitempty"`

	// Style is the declared style of the service.
	Style Id `json:"style"`
}
