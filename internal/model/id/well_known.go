// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id

// This file contains 'well-known' IDs, the idea being that we can avoid having to construct them at run-time.

const (
	// ArchFamilyC is the tag representing the C pseudo-architecture family.
	ArchFamilyC = "c"

	// ArchFamilyX86 is the tag representing the X86 architecture family.
	ArchFamilyX86 = "x86"
	// ArchFamilyArm is the tag representing the 32-bit Arm architecture family.
	ArchFamilyArm = "arm"
	// ArchFamilyAArch64 is the tag representing the 64-bit Arm architecture family.
	ArchFamilyAArch64 = "aarch64"
	// ArchFamilyPPC is the tag representing the PowerPC architecture family.
	ArchFamilyPPC = "ppc"

	// ArchVariantArm7 is the tag representing the arm7(-a) Arm variant.
	ArchVariantArm7 = "7"
	// ArchVariantArm8 is the tag representing the arm8(-a) Arm variant.
	ArchVariantArm8 = "8"
	// ArchVariantArmCortexA72 is the tag representing the Cortex-A72 Arm variant.
	// This variant is, for example, that used on the Raspberry Pi 4.
	ArchVariantArmCortexA72 = "cortex-a72"

	// ArchVariantAArch648 is the tag representing the arm8.x(-a) AArch64 variant.
	ArchVariantAArch648 = "8"

	// ArchSubVariantAArch6481 is the tag representing the arm8.1(-a) AArch64 variant.
	ArchSubVariantAArch6481 = "1"

	// ArchVariantPPC64LE is the tag representing the 64-bit little-endian PPC variant.
	ArchVariantPPC64LE = "64le"

	// ArchSubVariantPPCPOWER7 is the tag representing the POWER7 PPC sub-variant.
	ArchSubVariantPPCPOWER7 = "power7"
	// ArchSubVariantPPCPOWER8 is the tag representing the POWER8 PPC sub-variant.
	ArchSubVariantPPCPOWER8 = "power8"
	// ArchSubVariantPPCPOWER9 is the tag representing the POWER9 PPC sub-variant.
	ArchSubVariantPPCPOWER9 = "power9"

	// ArchVariantX8664 is the tag representing the 64-bit x86 variant.
	ArchVariantX8664 = "64"

	// ArchSubVariantX86Broadwell is the tag representing the Intel Broadwell x86-64 subvariant.
	ArchSubVariantX86Broadwell = "broadwell"
	// ArchSubVariantX86Skylake is the tag representing the Intel Skylake x86-64 subvariant.
	// This variant is, for example, that used in 2016 MacBook Pros.
	ArchSubVariantX86Skylake = "skylake"
)

var (
	// ArchC is the architecture ID for C.
	// (C isn't an architecture, but certain parts of the backend system are easier if we treat it as one.)
	ArchC = ID{[]string{ArchFamilyC}}

	// ArchX86 is the architecture ID for x86 (generic, assumed 32-bit).
	ArchX86 = ID{[]string{ArchFamilyX86}}
	// ArchX8664 is the architecture ID for x86-64.
	ArchX8664 = ID{[]string{ArchFamilyX86, ArchVariantX8664}}
	// ArchX86Broadwell is the architecture ID for x86-64 Broadwell.
	ArchX86Broadwell = ID{[]string{ArchFamilyX86, ArchVariantX8664, ArchSubVariantX86Broadwell}}
	// ArchX86Skylake is the architecture ID for x86-64 Skylake.
	ArchX86Skylake = ID{[]string{ArchFamilyX86, ArchVariantX8664, ArchSubVariantX86Skylake}}

	// ArchArm is the architecture ID for ARM (generic, 32-bit).
	ArchArm = ID{[]string{ArchFamilyArm}}
	// ArchArm7 is the architecture ID for arm7(-a).
	ArchArm7 = ID{[]string{ArchFamilyArm, ArchVariantArm7}}
	// ArchArm8 is the architecture ID for arm8(-a).
	ArchArm8 = ID{[]string{ArchFamilyArm, ArchVariantArm8}}
	// ArchArmCortexA72 is the architecture ID for arm Cortex-A72.
	ArchArmCortexA72 = ID{[]string{ArchFamilyArm, ArchVariantArmCortexA72}}

	// ArchAArch64 is the architecture ID for ARM (generic, 64-bit).
	ArchAArch64 = ID{[]string{ArchFamilyAArch64}}
	// ArchAArch648 is the architecture ID for ARM 64-bit version 8.x.
	ArchAArch648 = ID{[]string{ArchFamilyAArch64, ArchVariantAArch648}}
	// ArchAArch648 is the architecture ID for ARM 64-bit version 8.1.
	ArchAArch6481 = ID{[]string{ArchFamilyAArch64, ArchVariantAArch648, ArchSubVariantAArch6481}}

	// ArchPPC is the architecture ID for PowerPC.
	ArchPPC = ID{[]string{ArchFamilyPPC}}
	// ArchPPC64LE is the architecture ID for PowerPC64LE.
	ArchPPC64LE = ID{[]string{ArchFamilyPPC, ArchVariantPPC64LE}}
	// ArchPPCPOWER7 is the architecture ID for POWER7.
	ArchPPCPOWER7 = ID{[]string{ArchFamilyPPC, ArchVariantPPC64LE, ArchSubVariantPPCPOWER7}}
	// ArchPPCPOWER8 is the architecture ID for POWER8.
	ArchPPCPOWER8 = ID{[]string{ArchFamilyPPC, ArchVariantPPC64LE, ArchSubVariantPPCPOWER8}}
	// ArchPPCPOWER9 is the architecture ID for POWER9.
	ArchPPCPOWER9 = ID{[]string{ArchFamilyPPC, ArchVariantPPC64LE, ArchSubVariantPPCPOWER9}}

	// CStyleGCC is the compiler style ID for GCC.
	CStyleGCC = ID{[]string{"gcc"}}
)
