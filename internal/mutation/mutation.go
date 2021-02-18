// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mutation contains support for mutation testing using c4t.
package mutation

import (
	"fmt"
	"strconv"

	"github.com/1set/gut/ystring"
)

// Mutant is an identifier for a particular mutant.
//
// Since we only support a mutation testing setups with integer mutant identifiers, this is just uint64 for now.
type Mutant struct {
	// Name is the descriptive name of the mutant.
	Name Name

	// Index is the mutant index.
	//
	// The mutant index is what is passed into the mutation environment
	// variable, and is the basis for mutant definition by range.
	Index Index
}

// String gets a human-readable string representation of this mutant.
//
// If a name is available, the string will contain it.
func (m Mutant) String() string {
	if m.Name.IsZero() {
		return strconv.FormatUint(uint64(m.Index), 10)
	}
	return fmt.Sprintf("%s:%d", m.Name, m.Index)
}

// SetIndexIfZero sets this mutant's index to i if it is currently zero.
func (m *Mutant) SetIndexIfZero(i Index) {
	if m.Index == 0 {
		m.Index = i
	}
}

// NamedMutant creates a named mutant with index i, operator operator and variant variant.
// If operator is empty, the variant will not be recorded.
func NamedMutant(i Index, operator string, variant uint64) Mutant {
	m := AnonMutant(i)
	m.Name.Set(operator, variant)
	return m
}

// AnonMutant creates a mutant with index i, but no name.
func AnonMutant(i Index) Mutant {
	return Mutant{Index: i}
}

// Name is a human-readable name for mutants.
type Name struct {
	// Operator is the name of the mutant operator, if given.
	Operator string
	// Variant is the index of this particular mutant within its operator.
	Variant uint64
}

// IsZero gets whether this name appears to be the zero value.
func (n Name) IsZero() bool {
	return n.Operator == "" && n.Variant == 0
}

// String gets a string representation of this mutant name.
//
// The zero name returns the empty string; otherwise, the name is the operator name followed, if nonzero, by the
// variant number.
func (n Name) String() string {
	if n.Variant == 0 {
		return n.Operator
	}
	return fmt.Sprintf("%s%d", n.Operator, n.Variant)
}

// Set sets this name according to operator and variant.
// If operator is empty, we assume the mutant is unnamed, and clear the name to zero.
func (n *Name) Set(operator string, variant uint64) {
	n.Operator = operator
	// Zero the whole thing if the operator isn't given
	if ystring.IsEmpty(operator) {
		variant = 0
	}
	n.Variant = variant
}

// Index is the type of mutant indices.
type Index uint64

// EnvVar is the environment variable used for mutation testing.
//
// Some day, this might not be hard-coded.
const EnvVar = "C4_MUTATION"
