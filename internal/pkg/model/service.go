package model

// Service is a structure collecting information about a 'service' (a compiler or backend).
type Service struct {
	// ID is the (likely unqualified) ACT ID of the service.
	ID ID `toml:"id"`

	// IDQualified is true if Id is qualified by the machine ID.
	IDQualified bool `toml:"id_qualified,omitempty"`

	// MachineID is the ACT ID of the service's parent machine.
	// It may be empty if there is no need to track it.
	MachineID *ID `toml:"machine_id,omitempty"`

	// Style is the declared style of the service.
	Style ID `toml:"style"`
}

// FQID constructs the fully qualified ID of this service.
func (s Service) FQID() ID {
	if s.IDQualified {
		return s.ID
	}
	return ID{append(s.MachineID.tags, s.ID.tags...)}
}
