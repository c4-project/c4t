// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/compiler"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

var (
	ErrMalformedArchId    = errors.New("bad arch ID")
	ErrUnsupportedFamily  = errors.New("unsupported cpu family")
	ErrUnsupportedVariant = errors.New("unsupported cpu variant")
)

// mOptSet is a builder for GCC machine optimisation selections.
type mOptSet struct {
	// MArches contains the set of 'march' candidates.
	MArches stringhelp.Set
	// MCPUs contains the set of 'mcpu' candidates.
	MCPUs stringhelp.Set
	// AllowEmpty, if true, permits the selection of no-optimisation ("") rather than an march or a mcpu.
	AllowEmpty bool
}

func newMOptSet(allowEmpty bool) *mOptSet {
	return &mOptSet{
		MArches:    stringhelp.Set{},
		MCPUs:      stringhelp.Set{},
		AllowEmpty: allowEmpty,
	}
}

// AddArch adds arches to this set's march candidates.
func (m *mOptSet) AddArch(arches ...string) {
	m.MArches.Add(arches...)
}

// AddCPU adds cpus to this set's mcpu candidates.
func (m *mOptSet) AddCPU(cpus ...string) {
	m.MCPUs.Add(cpus...)
}

func (m *mOptSet) Strings() stringhelp.Set {
	narches := len(m.MArches)
	nstrs := narches + len(m.MCPUs)
	if m.AllowEmpty {
		nstrs++
	}
	nset := make(stringhelp.Set, nstrs)
	if m.AllowEmpty {
		nset.Add("")
	}
	for s := range m.MArches {
		nset.Add("arch=" + s)
	}
	for s := range m.MCPUs {
		nset.Add("cpu=" + s)
	}
	return nset
}

// DefaultMOpts adapts the GCC mopts calculation to the interface needed for a compiler.
func (g GCC) DefaultMOpts(c *compiler.Config) (stringhelp.Set, error) {
	return MOpts(c.Arch)
}

// MOpts gets the default 'm' invocations (march, mcpu, etc.) to consider for compilers with archID arch.
func MOpts(arch id.ID) (stringhelp.Set, error) {
	family, variant, ok := arch.Uncons()
	if !ok {
		return nil, fmt.Errorf("%w: empty", ErrMalformedArchId)
	}
	ms, err := mOptsFor(family, variant)
	if err != nil {
		return nil, err
	}
	return ms.Strings(), nil
}

func mOptsFor(family string, variant id.ID) (*mOptSet, error) {
	f, ok := invocationFamilies[family]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFamily, family)
	}
	return f(variant)
}

var invocationFamilies = map[string]func(id.ID) (*mOptSet, error){
	id.ArchFamilyArm: armMOpts,
	id.ArchFamilyPPC: ppcMOpts,
	id.ArchFamilyX86: x86MOpts,
}

func armMOpts(variant id.ID) (*mOptSet, error) {
	var err error
	// We disallow the empty variant to prevent issues whereby litmus7 generates dmb/dsb instructions, but gcc
	// selects a version of arm that doesn't understand them.
	set := newMOptSet(false)

	// Higher-numbered/more specific Arm variants append to lower-numbered/less specific variants,
	// hence the fallthrough-laden switch table.
	switch variant.String() {
	case id.ArchVariantArmCortexA72:
		set.MCPUs.Add("cortex-a72")
		fallthrough
	case id.ArchVariantArm8:
		set.MArches.Add("armv8-a")
		fallthrough
	case id.ArchVariantArm7:
		set.MArches.Add("armv7-a")
	case "":
		err = fmt.Errorf("%w: no variant (eg '7', '8', 'cortex-53') specified", ErrUnsupportedVariant)
	default:
		err = fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}

	return set, err
}

// ppcMOpts calculates the m-optimisation set for the PowerPC variant variant.
func ppcMOpts(variant id.ID) (*mOptSet, error) {
	mvar, svar, ok := variant.Uncons()
	if !ok {
		return nil, fmt.Errorf("%w: no variant (eg '64LE') specified", ErrUnsupportedVariant)
	}

	switch mvar {
	case id.ArchVariantPPC64LE:
		return ppc64LEMOpts(svar)
	default:
		return nil, fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}
}

// x8664MOpts calculates the m-optimisation set for the PowerPC64LE subvariant svar.
func ppc64LEMOpts(svar id.ID) (*mOptSet, error) {
	var err error
	set := newMOptSet(true)

	switch svar.String() {
	case id.ArchSubVariantPPCPOWER9:
		set.AddCPU("power9")
		fallthrough
	case id.ArchSubVariantPPCPOWER8:
		set.AddCPU("power8")
		fallthrough
	case id.ArchSubVariantPPCPOWER7:
		set.AddCPU("power7")
		fallthrough
	case "":
		set.AddCPU("powerpc64le", "native")
	default:
		return nil, fmt.Errorf("%w: unknown subvariant: %s", ErrUnsupportedVariant, svar)
	}

	return set, err
}

// x86MOpts calculates the m-optimisation set for the x86 variant variant.
func x86MOpts(variant id.ID) (*mOptSet, error) {
	mvar, svar, ok := variant.Uncons()
	if !ok {
		return nil, fmt.Errorf("%w: no variant (eg '64') specified", ErrUnsupportedVariant)
	}

	switch mvar {
	case id.ArchVariantX8664:
		return x8664MOpts(svar)
	default:
		return nil, fmt.Errorf("%w: unknown variant: %s", ErrUnsupportedVariant, variant)
	}
}

// x8664MOpts calculates the m-optimisation set for the x86-64 subvariant svar.
func x8664MOpts(svar id.ID) (*mOptSet, error) {
	var err error
	set := newMOptSet(true)

	switch svar.String() {
	case id.ArchSubVariantX86Skylake:
		set.AddArch("skylake")
		fallthrough
	case "":
		// TODO(@MattWindsor91): other subvariants?
		set.AddArch("x86_64", "native")
	default:
		return nil, fmt.Errorf("%w: unknown subvariant: %s", ErrUnsupportedVariant, svar)
	}

	return set, err
}
