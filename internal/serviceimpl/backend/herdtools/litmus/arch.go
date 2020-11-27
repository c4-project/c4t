// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"errors"
	"fmt"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/c4t/internal/model/id"
)

// archMap maps ACT architecture family/variant pairs to Litmus7 arch names.
// Each empty string mapping in a variant position is the 'default', or 'generic' architecture.
var archMap = map[string]map[string]string{
	id.ArchFamilyArm: {
		"": "ARM", // 32bit
	},
	id.ArchFamilyPPC: {
		"": "PPC",
	},
	id.ArchFamilyX86: {
		"":                  "X86", // 32-bit
		id.ArchVariantX8664: "X86_64",
	},
}

var (
	// ErrEmptyArch occurs when the arch ID sent to the Litmus backend is empty.
	ErrEmptyArch = errors.New("arch empty")
	// ErrBadArch occurs when the arch ID sent to the Litmus backend doesn't match any of the ones known to it.
	ErrBadArch = errors.New("arch family unknown")
)

func lookupArch(arch id.ID) (string, error) {
	f, v, _ := arch.Triple()
	if ystring.IsBlank(f) {
		return "", ErrEmptyArch
	}

	amap, ok := archMap[f]
	if !ok {
		mk, _ := id.MapKeys(archMap)
		return "", fmt.Errorf("%w: %s (valid: %q)", ErrBadArch, f, mk)
	}
	spec, ok := amap[v]
	if !ok {
		// Return the default if the variant doesn't have a specific match.
		return amap[""], nil
	}
	return spec, nil
}
