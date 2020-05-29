// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
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
	// No-operation.
	Nop Op = iota
	// Push all unconsumed inputs onto the stack.
	PushInputs
	// Push a specific input onto the stack, consuming it.
	// Takes a file argument.
	PushInput
	// Pop all of the inputs off the stack, compile them to an object, and push the name of the object onto the stack.
	CompileObj
	// Pop all of the inputs off the stack, compile them, and output the results to the output binary.
	CompileBin
	// Last is the last operation defined.
	Last = CompileBin
)

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
