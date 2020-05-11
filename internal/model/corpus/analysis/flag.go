// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

// flag programmatically represents, as a bitwise flag set, the possible classifications for a subject in a collation.
type flag int

const (
	// flagOk signifies the absence of collation flags.
	flagOk flag = 0
	// flagFlagged signifies that a subject was 'flagged'.
	flagFlagged flag = 1 << iota
	// flagCompileFail signifies a compile failure.
	flagCompileFail
	// flagCompileTimeout signifies a compile timeout.
	flagCompileTimeout
	// flagRunFailure signifies a runtime failure.
	flagRunFailure
	// flagRunTimeout signifies a runtime timeout.
	flagRunTimeout
)

// matches tests whether this flag matches expected.
// Generally this a bitwise test, except that flagOk only matches flagOk.
func (f flag) matches(expected flag) bool {
	if expected == flagOk {
		return f == flagOk
	}

	return (f & expected) == expected
}
