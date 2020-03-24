// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

// This file contains 'well-known' IDs, the idea being that we can avoid having to construct them at run-time.

var (
	// ArchX8664 is the ACT architecture ID for x86-64.
	ArchX8664 = ID{[]string{"x86", "64"}}

	// ArchArm is the ACT architecture ID for ARM.
	ArchArm = ID{[]string{"arm"}}

	// ArchPPC is the ACT architecture ID for PowerPC.
	ArchPPC = ID{[]string{"ppc"}}

	// CStyleGCC is the ACT compiler style for GCC.
	CStyleGCC = ID{[]string{"gcc"}}
)
