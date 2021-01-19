// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
)

// ppcMOpts calculates the m-optimisation set for the PowerPC variant variant.
func ppcMOpts(variant string, subvar id.ID) (*mOptSet, error) {
	switch variant {
	case id.ArchVariantPPC64LE:
		return ppc64LEMOpts(subvar)
	default:
		return nil, fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}
}

// ppc64LEMOpts calculates the m-optimisation set for the PowerPC64LE subvariant svar.
func ppc64LEMOpts(svar id.ID) (*mOptSet, error) {
	var err error
	set := newMOptSet(true)

	switch svar.String() {
	case id.ArchSubVariantPPCPOWER9:
		set.AddCPU("power9")
		fallthrough
	case id.ArchSubVariantPPCPOWER8:
		set.AddCPU("power8")
		fallthrough
	case id.ArchSubVariantPPCPOWER7:
		set.AddCPU("power7")
		fallthrough
	case "":
		set.AddCPU("powerpc64le", "native")
	default:
		return nil, fmt.Errorf("%w: unknown subvariant: %s", ErrUnsupportedVariant, svar)
	}

	return set, err
}
