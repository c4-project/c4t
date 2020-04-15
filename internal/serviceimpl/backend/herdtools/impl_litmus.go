// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
)

// Litmus describes the parts of a Litmus invocation that are specific to Herd.
type Litmus struct{}

// ParseStateCount parses a Litmus state count.
func (l Litmus) ParseStateCount(fields []string) (uint64, error) {
	if nf := len(fields); nf != 3 {
		return 0, fmt.Errorf("%w: expected three fields, got %d", ErrBadStateCount, nf)
	}
	if f := fields[0]; f != "Histogram" {
		return 0, fmt.Errorf("%w: expected first word to be 'Histogram', got %q", ErrBadStateCount, f)
	}
	if f := fields[2]; f != "states)" {
		return 0, fmt.Errorf("%w: expected last word to be 'states)', got %q", ErrBadStateCount, f)
	}
	return strconv.ParseUint(strings.TrimPrefix(fields[1], "("), 10, 64)
}

func (l Litmus) ParseStateLine(tt TestType, fields []string) (*StateLine, error) {
	nf := len(fields)
	if nf == 0 {
		return nil, fmt.Errorf("%w: expected at least one field", ErrBadStateLine)
	}

	var (
		s   StateLine
		err error
	)

	// The start of a Litmus state line is always of the form N:>x=y;, where:
	// - N is the number of times the state was observed;
	// - : is * when the line is 'unusual' (a witness for an 'allowed' test, or a counter for a 'required' test);
	// - x=y; is the first mapping in the state (with no space between it and the >).
	//
	// There may be some space after N, which means we can't rely on the field split.
	line := strings.Join(fields, " ")

	errf := func(why string) (*StateLine, error) {
		return nil, fmt.Errorf("%w: %s, got %q", ErrBadStateLine, why, line)
	}

	splits := strings.Split(line, ">")
	if len(splits) != 2 {
		return errf("expected exactly one '>'")
	}

	var meta string
	meta, s.Rest = splits[0], strings.Fields(splits[1])

	lfm := len(meta)
	if lfm == 0 {
		return errf("expected metadata before '>'")
	}

	if s.NOccurs, err = parseNOccurs(meta[:lfm-1]); err != nil {
		return nil, err
	}

	s.Tag, err = parseTagSigil(tt, rune(meta[lfm-1]))
	return &s, err
}

func parseNOccurs(raw string) (uint64, error) {
	nOccursStr := strings.TrimSpace(raw)
	return strconv.ParseUint(nOccursStr, 10, 64)
}

const (
	// sigilEmph appears on witnessing status lines in an 'allowed' test, and counter status lines in a 'required' test.
	sigilEmph = '*'
	// sigilNorm appears on counter status lines in an 'allowed' test, and witness status lines in a 'required' test.
	sigilNorm = ':'
)

func parseTagSigil(tt TestType, sigil rune) (obs.Tag, error) {
	switch tt {
	case TTAllowed:
		return parseTagSigilLine(sigil, obs.TagWitness, obs.TagCounter)
	case TTRequired:
		return parseTagSigilLine(sigil, obs.TagCounter, obs.TagWitness)
	default:
		return obs.TagUnknown, fmt.Errorf("%w: unknown test type %v", ErrBadStateLine, tt)
	}
}

func parseTagSigilLine(sigil rune, onEmph, onNorm obs.Tag) (obs.Tag, error) {
	switch sigil {
	case sigilEmph:
		return onEmph, nil
	case sigilNorm:
		return onNorm, nil
	default:
		return obs.TagUnknown, fmt.Errorf("%w: unknown sigil %q", ErrBadStateLine, sigil)
	}
}

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
func (l Litmus) Args(j job.Harness, r service.RunInfo) ([]string, error) {
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
