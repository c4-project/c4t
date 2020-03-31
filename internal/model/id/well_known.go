// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id

// This file contains 'well-known' IDs, the idea being that we can avoid having to construct them at run-time.

const (
	// ArchFamilyX86 is the tag representing the X86 architecture family.
	ArchFamilyX86 = "x86"
	// ArchFamilyArm is the tag representing the 32-bit Arm architecture family.
	ArchFamilyArm = "arm"
	// ArchFamilyPPC is the tag representing the PowerPC architecture family.
	ArchFamilyPPC = "ppc"
)

var (
	// ArchX86 is the ACT architecture ID for x86 (generic, assumed 32-bit).
	ArchX86 = ID{[]string{ArchFamilyX86}}

	// ArchX8664 is the ACT architecture ID for x86-64.
	ArchX8664 = ID{[]string{ArchFamilyX86, "64"}}

	// ArchArm is the ACT architecture ID for ARM (32-bit).
	ArchArm = ID{[]string{ArchFamilyArm}}

	// ArchPPC is the ACT architecture ID for PowerPC.
	ArchPPC = ID{[]string{ArchFamilyPPC}}

	// CStyleGCC is the ACT compiler style for GCC.
	CStyleGCC = ID{[]string{"gcc"}}
)
