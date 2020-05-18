// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis

import (
	"github.com/MattWindsor91/act-tester/internal/model/status"
)

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

	// FlagFail is the union of all failure flags.
	FlagFail = FlagCompileFail | FlagRunFail
	// FlagTimeout is the union of all timeout flags.
	FlagTimeout = FlagCompileTimeout | FlagRunTimeout
)

// Matches tests whether this Flag Matches expected.
// Generally this a bitwise test, except that FlagOk only Matches FlagOk.
func (f Flag) Matches(expected Flag) bool {
	if expected == FlagOk {
		return f == FlagOk
	}

	return (f & expected) == expected
}

// statusFlags Matches statuses to flags.
var statusFlags = [status.Num]Flag{
	status.Ok:             FlagOk,
	status.Flagged:        FlagFlagged,
	status.CompileTimeout: FlagCompileTimeout,
	status.CompileFail:    FlagCompileFail,
	status.RunTimeout:     FlagRunTimeout,
	status.RunFail:        FlagRunFail,
}
