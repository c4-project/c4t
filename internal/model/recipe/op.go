// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Op is the enumeration of kinds of operation that can be in a recipe.
type Op uint8

const (
	// Nop is a no-operation.
	Nop Op = iota
	// PushInputs is a compile instruction that pushes all unconsumed inputs (matching the filekind argument, if given).
	PushInputs
	// PushInput is a compile instruction that pushes a specific input onto the stack, consuming it.
	// Takes a file argument.
	PushInput
	// CompileObj pops inputs off the stack, compile them to an object, and push the name of the object onto the stack.
	CompileObj
	// CompileObj pops inputs off the stack, compile them, and output the results to the output binary.
	CompileExe

	// Last is the last operation defined.
	Last = CompileExe
)

//go:generate stringer -type Op

// OpFromString tries to convert a string into an Op.
func OpFromString(s string) (Op, error) {
	for i := Nop; i <= Last; i++ {
		if strings.EqualFold(s, i.String()) {
			return i, nil
		}
	}
	return Nop, fmt.Errorf("unknown Op: %q", s)
}

// MarshalJSON marshals an op to JSON using its string form.
func (i Op) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON unmarshals an op from JSON using its string form.
func (i *Op) UnmarshalJSON(bytes []byte) error {
	var (
		is  string
		err error
	)
	if err = json.Unmarshal(bytes, &is); err != nil {
		return err
	}
	*i, err = OpFromString(is)
	return err
}
