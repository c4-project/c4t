// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

// TODO(@MattWindsor91): backends are unusually complex, mostly because I originally expected that backends would need
// to be accessible from machines other than the test compiler.  It's turned out that this isn't actually the case, and
// so there is scope for scaling down the complexity eventually.

// Backend collects the test-relevant information about a backend.
type Backend struct {
	// CompilerID is the (likely unqualified) ACT CompilerID of the backend.
	ID ID `toml:"id"`

	// IDQualified is true if Id is qualified by the machine CompilerID.
	IDQualified bool `toml:"id_qualified,omitempty"`

	// MachineID is the ACT CompilerID of the backend's parent machine.
	// It may be empty if there is no need to track it.
	MachineID *ID `toml:"machine_id,omitempty"`

	// Style is the declared style of the service.
	Style ID `toml:"style"`
}

// FQID constructs the fully qualified CompilerID of this service.
func (b Backend) FQID() ID {
	if b.IDQualified || b.MachineID == nil {
		return b.ID
	}
	return b.MachineID.Join(b.ID)
}
