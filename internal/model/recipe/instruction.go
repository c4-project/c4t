// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import "fmt"

//go:generate stringer -type Op

// Instruction represents a single instruction in a recipe.
//
// Instructions target a stack machine in the machine node.
type Instruction struct {
	// Op is the opcode.
	Op Op `json:"op"`

	// File is, if applicable, the file argument to the instruction.
	File string `json:"file,omitempty"`
}

// String produces a human-readable string representation of this instruction.
func (i Instruction) String() string {
	switch i.Op {
	case PushInput:
		return fmt.Sprintf("%s %q", i.Op.String(), i.File)
	default:
		return i.Op.String()
	}
}

// CompileBinInst produces a 'compile binary' instruction.
func CompileBinInst() Instruction {
	return Instruction{Op: CompileBin}
}

// PushInputInst produces a 'push input' instruction.
func PushInputInst(file string) Instruction {
	return Instruction{Op: PushInput, File: file}
}
