// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"strconv"
	"strings"
)

const (
	// MutantHitPrefix is the prefix of lines from compilers specifying that a mutant has been hit.
	MutantHitPrefix = "MUTATION HIT:"
	// MutantSelectPrefix is the prefix of lines from compilers specifying that a mutant has been selected.
	MutantSelectPrefix = "MUTATION SELECTED:"
)

// ScanLine scans line for mutant hit and selection hints, and calls the appropriate callback.
func ScanLine(line string, onHit, onSelect func(uint64)) {
	line = strings.TrimSpace(line)

	for prefix, f := range map[string]func(uint64){
		MutantHitPrefix:    onHit,
		MutantSelectPrefix: onSelect,
	} {
		if strings.HasPrefix(line, prefix) {
			scanLineAfterPrefix(strings.TrimPrefix(line, prefix), f)
		}
	}
}

func scanLineAfterPrefix(line string, f func(uint64)) {
	// Some of the lines contain things other than the mutant number.
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}

	n, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return
	}

	f(n)
}
