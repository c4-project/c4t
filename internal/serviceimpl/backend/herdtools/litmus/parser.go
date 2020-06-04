// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/model/obs"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/parser"
)

// ParseStateCount parses a Litmus state count.
func (l Litmus) ParseStateCount(fields []string) (uint64, error) {
	if nf := len(fields); nf != 3 {
		return 0, fmt.Errorf("%w: expected three fields, got %d", parser.ErrBadStateCount, nf)
	}
	if f := fields[0]; f != "Histogram" {
		return 0, fmt.Errorf("%w: expected first word to be 'Histogram', got %q", parser.ErrBadStateCount, f)
	}
	if f := fields[2]; f != "states)" {
		return 0, fmt.Errorf("%w: expected last word to be 'states)', got %q", parser.ErrBadStateCount, f)
	}
	return strconv.ParseUint(strings.TrimPrefix(fields[1], "("), 10, 64)
}

func (l Litmus) ParseStateLine(tt parser.TestType, fields []string) (*parser.StateLine, error) {
	nf := len(fields)
	if nf == 0 {
		return nil, fmt.Errorf("%w: expected at least one field", parser.ErrBadStateLine)
	}

	// The start of a Litmus state line is always of the form N:>x=y;, where:
	// - N is the number of times the state was observed;
	// - : is * when the line is 'unusual' (a witness for an 'allowed' test, or a counter for a 'required' test);
	// - x=y; is the first mapping in the state (with no space between it and the >).
	//
	// There may be some space after N, which means we can't rely on the field split.
	line := parseLine{line: strings.Join(fields, " "), tt: tt}
	return line.parse()
}

// parseLine is an intermediate struct used for parsing a state line.
type parseLine struct {
	line string
	tt   parser.TestType
}

func (l *parseLine) parse() (*parser.StateLine, error) {
	splits := strings.Split(l.line, ">")
	if len(splits) != 2 {
		return l.errorOutf("expected exactly one '>'")
	}

	return l.parseWithMeta(splits[0], strings.Fields(splits[1]))
}

func (l *parseLine) parseWithMeta(meta string, rest []string) (*parser.StateLine, error) {
	var (
		s   parser.StateLine
		err error
	)

	lfm := len(meta)
	if lfm == 0 {
		return l.errorOutf("expected metadata before '>'")
	}

	if s.NOccurs, err = parseNOccurs(meta[:lfm-1]); err != nil {
		return nil, err
	}
	s.Rest = rest
	s.Tag, err = l.parseTagSigil(rune(meta[lfm-1]))
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

func (l *parseLine) parseTagSigil(sigil rune) (obs.Tag, error) {
	switch l.tt {
	case parser.Allowed:
		return parseTagSigilLine(sigil, obs.TagWitness, obs.TagCounter)
	case parser.Required:
		return parseTagSigilLine(sigil, obs.TagCounter, obs.TagWitness)
	default:
		return obs.TagUnknown, l.errorf("unknown test type %v", l.tt)
	}
}

func (l *parseLine) errorf(format string, arg ...interface{}) error {
	why := fmt.Sprintf(format, arg...)
	return fmt.Errorf("%w: %s, got %q", parser.ErrBadStateLine, why, l.line)
}

func (l *parseLine) errorOutf(format string, arg ...interface{}) (*parser.StateLine, error) {
	return nil, l.errorf(format, arg...)
}

func parseTagSigilLine(sigil rune, onEmph, onNorm obs.Tag) (obs.Tag, error) {
	switch sigil {
	case sigilEmph:
		return onEmph, nil
	case sigilNorm:
		return onNorm, nil
	default:
		return obs.TagUnknown, fmt.Errorf("%w: unknown sigil %q", parser.ErrBadStateLine, sigil)
	}
}
