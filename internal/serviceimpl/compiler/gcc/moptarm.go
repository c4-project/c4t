// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
)

func armMOpts(variant string, _ id.ID) (*mOptSet, error) {
	var err error
	// We disallow the empty variant to prevent issues whereby litmus7 generates dmb/dsb instructions, but gcc
	// selects a version of arm that doesn't understand them.
	set := newMOptSet(false)

	// Higher-numbered/more specific Arm variants append to lower-numbered/less specific variants,
	// hence the fallthrough-laden switch table.
	switch variant {
	case id.ArchVariantArmCortexA72:
		set.MCPUs.Add("cortex-a72")
		fallthrough
	case id.ArchVariantArm8:
		set.MArches.Add("armv8-a")
		fallthrough
	case id.ArchVariantArm7:
		set.MArches.Add("armv7-a")
	case "":
		err = fmt.Errorf("%w: no variant (eg '7', '8', 'cortex-53') specified", ErrUnsupportedVariant)
	default:
		err = fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}

	return set, err
}
