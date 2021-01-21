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

// ScanLines scans each line in lines, building a map of mutant numbers to hit counts.
// If a mutant is present in the map, it was selected, even if its hit count is 0.
func ScanLines(lines []string) map[uint64]uint64 {
	mp := make(map[uint64]uint64)
	onHit := func(i uint64) {
		mp[i] = mp[i] + 1
	}
	onSelect := func(i uint64) {
		// Defines mp[i] with 0 if it hasn't already been defined.
		mp[i] = mp[i] + 0
	}
	for _, l := range lines {
		ScanLine(l, onHit, onSelect)
	}
	return mp
}

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
