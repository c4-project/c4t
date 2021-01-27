// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

const (
	// MutantHitPrefix is the prefix of lines from compilers specifying that a mutant has been hit.
	MutantHitPrefix = "MUTATION HIT:"
	// MutantSelectPrefix is the prefix of lines from compilers specifying that a mutant has been selected.
	MutantSelectPrefix = "MUTATION SELECTED:"
)

// ScanLines scans each line in r, building a map of mutant numbers to hit counts.
// If a mutant is present in the map, it was selected, even if its hit count is 0.
func ScanLines(r io.Reader) map[Mutant]uint64 {
	sc := bufio.NewScanner(r)

	mp := make(map[Mutant]uint64)
	onHit := func(i Mutant) {
		mp[i] = mp[i] + 1
	}
	onSelect := func(i Mutant) {
		// Defines mp[i] with 0 if it hasn't already been defined.
		mp[i] = mp[i] + 0
	}
	for sc.Scan() {
		ScanLine(sc.Text(), onHit, onSelect)
	}
	return mp
}

// ScanLine scans line for mutant hit and selection hints, and calls the appropriate callback.
func ScanLine(line string, onHit, onSelect func(Mutant)) {
	line = strings.TrimSpace(line)

	for prefix, f := range map[string]func(Mutant){
		MutantHitPrefix:    onHit,
		MutantSelectPrefix: onSelect,
	} {
		if strings.HasPrefix(line, prefix) {
			scanLineAfterPrefix(strings.TrimPrefix(line, prefix), f)
		}
	}
}

func scanLineAfterPrefix(line string, f func(Mutant)) {
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
