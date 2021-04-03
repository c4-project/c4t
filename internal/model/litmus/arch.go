// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/c4-project/c4t/internal/helper/errhelp"

	"github.com/c4-project/c4t/internal/id"
)

// archMap maps C4 architecture family/variant pairs to Litmus7 arch names.
// Each empty string mapping in a variant position is the 'default', or 'generic' architecture.
// The empty ID doesn't map to an architecture.
var archMap = map[string]map[string]string{
	id.ArchFamilyC: {
		"": "C",
	},
	id.ArchFamilyArm: {
		"": "ARM", // 32-bit
	},
	id.ArchFamilyAArch64: {
		"": "AArch64", // 64-bit
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

// ArchOfLitmus tries to look up the C4 identifier of a Litmus architecture.
func ArchOfLitmus(arch string) (id.ID, error) {
	for fam, vs := range archMap {
		for avar, arch2 := range vs {
			if strings.EqualFold(arch, arch2) {
				if avar == "" {
					return id.New(fam)
				}
				return id.New(fam, avar)
			}
		}
	}
	return id.ID{}, fmt.Errorf("%w: unsupported Litmus architecture %q", ErrBadArch, arch)
}

// ArchToLitmus tries to look up the Litmus name of a C4 architecture ID.
func ArchToLitmus(arch id.ID) (string, error) {
	if arch.IsEmpty() {
		return "", ErrEmptyArch
	}

	f, v, _ := arch.Triple()
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

// ArchOfFile tries to divine the architecture ID of a Litmus test by reading its first line from file fpath.
func ArchOfFile(fpath string) (id.ID, error) {
	r, err := os.Open(fpath)
	if err != nil {
		return id.ID{}, err
	}
	br := bufio.NewReader(r)
	hdr, rerr := br.ReadString('\n')
	cerr := r.Close()
	if err := errhelp.FirstError(rerr, cerr); err != nil {
		return id.ID{}, err
	}

	hfields := strings.Fields(hdr)
	if len(hfields) == 0 {
		return id.ID{}, fmt.Errorf("%w: litmus file has empty initial line", ErrEmptyArch)
	}
	return ArchOfLitmus(hfields[0])
}

// PopulateArchFromFile sets this litmus test's architecture to that from calling ArchOfFile over its defined Filepath.
func (l *Litmus) PopulateArchFromFile() error {
	var err error
	l.Arch, err = ArchOfFile(l.Filepath())
	return err
}
