// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe

import (
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/model/filekind"
)

// Option is a functional option for a recipe.
type Option func(*Recipe) error

// ErrNotTakingInstructions occurs if we try to add instructions to a recipe where they don't make sense.
var ErrNotTakingInstructions = errors.New("this output type doesn't support instructions")

// Options applies multiple options to a recipe.
func Options(os ...Option) Option {
	return func(r *Recipe) error {
		for _, o := range os {
			if err := o(r); err != nil {
				return err
			}
		}
		return nil
	}
}

// AddFiles adds each file in fs to the recipe.
func AddFiles(fs ...string) Option {
	return func(r *Recipe) error {
		r.Files = append(r.Files, fs...)
		return nil
	}
}

// AddInstructions adds each instruction in ins to the recipe.
func AddInstructions(ins ...Instruction) Option {
	return func(r *Recipe) error {
		if r.Output == OutNothing {
			return fmt.Errorf("%w: %s", ErrNotTakingInstructions, r.Output)
		}
		r.Instructions = append(r.Instructions, ins...)
		return nil
	}
}

// CompileFileToObj adds a set of instructions that compile the named C input to an object file.
func CompileFileToObj(file string) Option {
	return AddInstructions(
		PushInputInst(file),
		CompileObjInst(1),
	)
}

// CompileAllCToExe adds a set of instructions that compile all C inputs to an executable at path.
func CompileAllCToExe() Option {
	return Options(
		AddInstructions(
			PushInputsInst(filekind.CSrc),
			CompileExeInst(PopAll),
		),
	)
}
