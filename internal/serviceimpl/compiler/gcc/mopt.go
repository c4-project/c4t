// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"errors"
	"fmt"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/model/id"
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
func (g GCC) DefaultMOpts(c *compiler.Compiler) (stringhelp.Set, error) {
	return MOpts(c.Arch)
}

// MOpts gets the default 'm' invocations (march, mcpu, etc.) to consider for compilers with archID arch.
func MOpts(arch id.ID) (stringhelp.Set, error) {
	family, variant, subvar := arch.Triple()
	if ystring.IsEmpty(family) {
		return nil, fmt.Errorf("%w: empty", ErrMalformedArchId)
	}
	ms, err := mOptsFor(family, variant, subvar)
	if err != nil {
		return nil, err
	}
	return ms.Strings(), nil
}

func mOptsFor(family, variant string, subvar id.ID) (*mOptSet, error) {
	f, ok := invocationFamilies[family]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFamily, family)
	}
	return f(variant, subvar)
}

var invocationFamilies = map[string]func(string, id.ID) (*mOptSet, error){
	id.ArchFamilyAArch64: aarch64MOpts,
	id.ArchFamilyArm:     armMOpts,
	id.ArchFamilyPPC:     ppcMOpts,
	id.ArchFamilyX86:     x86MOpts,
}
