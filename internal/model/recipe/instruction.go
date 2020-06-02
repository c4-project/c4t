// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
)

// Instruction represents a single instruction in a recipe.
//
// Instructions target a stack machine in the machine node.
type Instruction struct {
	// Op is the opcode.
	Op Op `json:"op"`

	// File is, if applicable, the file argument to the instruction.
	File string `json:"file,omitempty"`

	// FileKind is, if applicable, the file kind argument to the instruction.
	FileKind filekind.Kind `json:"file_kind,omitempty"`
}

// String produces a human-readable string representation of this instruction.
func (i Instruction) String() string {
	switch i.Op {
	case PushInput:
		return fmt.Sprintf("%s %q", i.Op.String(), i.File)
	case PushInputs:
		return fmt.Sprintf("%s %q", i.Op.String(), i.FileKind)
	default:
		return i.Op.String()
	}
}

// CompileExeInst produces a 'compile binary' instruction.
func CompileExeInst() Instruction {
	return Instruction{Op: CompileExe}
}

// CompileObjInst produces a 'compile object' instruction.
func CompileObjInst() Instruction {
	return Instruction{Op: CompileObj}
}

// PushInputInst produces a 'push input' instruction.
func PushInputInst(file string) Instruction {
	return Instruction{Op: PushInput, File: file}
}

// PushInputsInst produces a 'push inputs' instruction.
func PushInputsInst(kind filekind.Kind) Instruction {
	return Instruction{Op: PushInputs, FileKind: kind}
}
