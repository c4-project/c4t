// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package status

// Flag programmatically represents, as a bitwise Flag set, the possible classifications for a subject.
type Flag int

const (
	// FlagFiltered signifies that a subject was filtered out.
	FlagFiltered Flag = 1 << iota
	// FlagFlagged signifies that a subject was 'flagged'.
	FlagFlagged
	// FlagCompileFail signifies a compile failure.
	FlagCompileFail
	// FlagCompileTimeout signifies a compile timeout.
	FlagCompileTimeout
	// FlagRunFail signifies a runtime failure.
	FlagRunFail
	// FlagRunTimeout signifies a runtime timeout.
	FlagRunTimeout

	// FlagFail is the union of all failure flags.
	FlagFail = FlagCompileFail | FlagRunFail
	// FlagTimeout is the union of all timeout flags.
	FlagTimeout = FlagCompileTimeout | FlagRunTimeout
	// FlagBad is the union of all 'bad' flags; it should match the calculation in Status.IsBad.
	FlagBad = FlagFail | FlagTimeout | FlagFlagged
)

// Matches tests whether this Flag has all flag bits in expected present.
func (f Flag) Matches(expected Flag) bool {
	return (f & expected) == expected
}

// MatchesStatus tests whether this Flag matches the expected status.
// Generally, this is a Matches test for Status.Flag, except that Ok only matches an absence of other flags.
func (f Flag) MatchesStatus(expected Status) bool {
	if expected == Ok {
		return f == 0
	}
	return f.Matches(expected.Flag())
}

// statusFlags matches statuses to flags.
var statusFlags = [Last + 1]Flag{
	// Unknown and Ok don't have any flags set.
	Filtered:       FlagFiltered,
	Flagged:        FlagFlagged,
	CompileTimeout: FlagCompileTimeout,
	CompileFail:    FlagCompileFail,
	RunTimeout:     FlagRunTimeout,
	RunFail:        FlagRunFail,
}

// Flag gets the flag equivalent of this status.
func (i Status) Flag() Flag {
	return statusFlags[i]
}
