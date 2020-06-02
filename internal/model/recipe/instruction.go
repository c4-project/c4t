// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import (
	"fmt"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
)

// PopAll is the value to pass to NPops to ask the instruction to pop all applicable files off the stack.
const PopAll = 0

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

	// npops is, if applicable and nonzero, the maximum number of items to pop off the file stack.
	NPops int `json:"npops,omitempty"`
}

// String produces a human-readable string representation of this instruction.
func (i Instruction) String() string {
	switch i.Op {
	case CompileExe:
		fallthrough
	case CompileObj:
		return fmt.Sprintf("%s %s", i.Op, npopString(i.NPops))
	case PushInput:
		return fmt.Sprintf("%s %q", i.Op, i.File)
	case PushInputs:
		return fmt.Sprintf("%s %q", i.Op, i.FileKind)
	default:
		return i.Op.String()
	}
}

// npopString returns 'ALL' if npops requests popping all files, or npops as a string otherwise.
func npopString(npops int) string {
	if npops <= PopAll {
		return "ALL"
	}
	return strconv.Itoa(npops)
}

// CompileExeInst produces a 'compile binary' instruction.
func CompileExeInst(npops int) Instruction {
	return Instruction{Op: CompileExe, NPops: npops}
}

// CompileObjInst produces a 'compile object' instruction.
func CompileObjInst(npops int) Instruction {
	return Instruction{Op: CompileObj, NPops: npops}
}

// PushInputInst produces a 'push input' instruction.
func PushInputInst(file string) Instruction {
	return Instruction{Op: PushInput, File: file}
}

// PushInputsInst produces a 'push inputs' instruction.
func PushInputsInst(kind filekind.Kind) Instruction {
	return Instruction{Op: PushInputs, FileKind: kind}
}
