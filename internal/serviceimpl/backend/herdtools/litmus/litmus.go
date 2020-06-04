// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains the parts of a Herdtools backend specific to herd7.
package litmus

import (
	"errors"
	"fmt"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Litmus describes the parts of a Litmus invocation that are specific to Herd.
type Litmus struct{}

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

// Args deduces the appropriate arguments for running Litmus on job j, with the merged run information r.
func (l Litmus) Args(j job.Lifter, r service.RunInfo) ([]string, error) {
	larch, err := lookupArch(j.Arch)
	if err != nil {
		return nil, fmt.Errorf("when looking up -carch: %w", err)
	}
	args := []string{
		"-o", j.OutDir,
		"-carch", larch,
		"-c11", "true",
	}
	args = append(args, r.Args...)
	args = append(args, j.InFile)
	return args, nil
}
