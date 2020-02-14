package model

// TODO(@MattWindsor91): backends are unusually complex, mostly because I originally expected that backends would need
// to be accessible from machines other than the test runner.  It's turned out that this isn't actually the case, and
// so there is scope for scaling down the complexity eventually.

// Backend collects the test-relevant information about a backend.
type Backend struct {
	// ID is the (likely unqualified) ACT ID of the backend.
	ID ID `toml:"id"`

	// IDQualified is true if Id is qualified by the machine ID.
	IDQualified bool `toml:"id_qualified,omitempty"`

	// MachineID is the ACT ID of the backend's parent machine.
	// It may be empty if there is no need to track it.
	MachineID *ID `toml:"machine_id,omitempty"`

	// Style is the declared style of the service.
	Style ID `toml:"style"`
}

// FQID constructs the fully qualified ID of this service.
func (b Backend) FQID() ID {
	if b.IDQualified {
		return b.ID
	}
	return ID{append(b.MachineID.tags, b.ID.tags...)}
}
