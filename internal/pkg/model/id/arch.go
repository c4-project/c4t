// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

var (
	// ArchX8664 is the ACT architecture CompilerID for x86-64.
	ArchX8664 = ID{[]string{"x86", "64"}}

	// ArchArm is the ACT architecture CompilerID for ARM.
	ArchArm = ID{[]string{"arm"}}

	// ArchPPC is the ACT architecture CompilerID for PowerPC.
	ArchPPC = ID{[]string{"ppc"}}
)
