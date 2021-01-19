// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
)

// x86MOpts calculates the m-optimisation set for the triple (x86, variant, subvar).
func x86MOpts(variant string, subvar id.ID) (*mOptSet, error) {
	switch variant {
	case id.ArchVariantX8664:
		return x8664MOpts(subvar)
	default:
		return nil, fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}
}

// x8664MOpts calculates the m-optimisation set for the x86-64 subvariant svar.
func x8664MOpts(svar id.ID) (*mOptSet, error) {
	var err error
	set := newMOptSet(true)

	switch svar.String() {
	case id.ArchSubVariantX86Skylake:
		set.AddArch("skylake")
		// TODO(@MattWindsor91): is this fallthrough safe?
		fallthrough
	case id.ArchSubVariantX86Broadwell:
		set.AddArch("broadwell")
		fallthrough
	case "":
		// TODO(@MattWindsor91): other subvariants?
		set.AddArch("x86-64", "native")
	default:
		return nil, fmt.Errorf("%w: unknown subvariant: %s", ErrUnsupportedVariant, svar)
	}

	return set, err
}
