// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/filekind"

	"github.com/c4-project/c4t/internal/model/recipe"
)

// ExampleCompileExeInst is a runnable example for CompileExeInst.
func ExampleCompileExeInst() {
	fmt.Println(recipe.CompileExeInst(recipe.PopAll))
	fmt.Println(recipe.CompileExeInst(1))

	// Output:
	// CompileExe ALL
	// CompileExe 1
}

// ExampleCompileObjInst is a runnable example for CompileObjInst.
func ExampleCompileObjInst() {
	fmt.Println(recipe.CompileObjInst(recipe.PopAll))
	fmt.Println(recipe.CompileObjInst(1))

	// Output:
	// CompileObj ALL
	// CompileObj 1
}

// ExamplePushInputInst is a runnable example for PushInputInst.
func ExamplePushInputInst() {
	fmt.Println(recipe.PushInputInst("foo.c"))

	// Output:
	// PushInput "foo.c"
}

// ExamplePushInputsInst is a runnable example for PushInputsInst.
func ExamplePushInputsInst() {
	fmt.Println(recipe.PushInputsInst(filekind.C))

	// Output:
	// PushInputs c
}
