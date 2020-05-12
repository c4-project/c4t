// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

// Flag programmatically represents, as a bitwise Flag set, the possible classifications for a subject in a collation.
type Flag int

const (
	// FlagOk signifies the absence of collation flags.
	FlagOk Flag = 0
	// FlagFlagged signifies that a subject was 'flagged'.
	FlagFlagged Flag = 1 << iota
	// FlagCompileFail signifies a compile failure.
	FlagCompileFail
	// FlagCompileTimeout signifies a compile timeout.
	FlagCompileTimeout
	// FlagRunFail signifies a runtime failure.
	FlagRunFail
	// FlagRunTimeout signifies a runtime timeout.
	FlagRunTimeout
)

// matches tests whether this Flag matches expected.
// Generally this a bitwise test, except that FlagOk only matches FlagOk.
func (f Flag) matches(expected Flag) bool {
	if expected == FlagOk {
		return f == FlagOk
	}

	return (f & expected) == expected
}
