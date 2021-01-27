// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
)

func aarch64MOpts(variant string, subvar id.ID) (*mOptSet, error) {
	// Unlike 32-bit, on Arm we should be safe to emit for native/generic.

	var err error
	set := newMOptSet(true)

	switch variant {
	// TODO(@MattWindsor91): Apple M1
	case id.ArchVariantArm8:
		set.MArches.Add("armv8-a")
		err = aarch648MOpts(set, subvar.String())
		fallthrough
	case "":
		set.AddCPU("generic", "native")
	default:
		err = fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}

	return set, err
}

// aarch648MOpts adds to set mopts useful for subvariants of 64-bit armv8.
func aarch648MOpts(set *mOptSet, subvar string) error {
	switch subvar {
	// TODO(@MattWindsor91): Apple M1
	case id.ArchSubVariantAArch6481:
		set.MArches.Add("armv8.1-a")
		fallthrough
	case "":
		return nil
	default:
		return fmt.Errorf("%w: unknown subvariant: %s", ErrUnsupportedVariant, subvar)
	}
}
