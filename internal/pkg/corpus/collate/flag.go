// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate

type collationFlag int

const (
	ccOk      collationFlag = 0
	ccCompile collationFlag = 1 << iota
	ccFlag
	ccRunFailure
	ccTimeout
)

// matches tests whether this flag matches expected.
// Generally this a bitwise test, except that ccOk only matches ccOk.
func (f collationFlag) matches(expected collationFlag) bool {
	if expected == ccOk {
		return f == ccOk
	}

	return (f & expected) == expected
}
