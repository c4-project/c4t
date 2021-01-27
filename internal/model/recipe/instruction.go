// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import (
	"strconv"
	"strings"

	"github.com/c4-project/c4t/internal/model/filekind"
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

	// NPops is, if applicable and nonzero, the maximum number of items to pop off the file stack.
	NPops int `json:"npops,omitempty"`
}

// String produces a human-readable string representation of this instruction.
func (i Instruction) String() string {
	strs := []string{i.Op.String()}

	switch i.Op {
	case CompileExe:
		fallthrough
	case CompileObj:
		strs = append(strs, npopString(i.NPops))
	case PushInput:
		strs = append(strs, strconv.Quote(i.File))
	case PushInputs:
		strs = append(strs, i.FileKind.String())
	}

	return strings.Join(strs, " ")
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
