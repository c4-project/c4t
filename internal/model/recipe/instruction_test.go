// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// ExampleCompileBinInst is a runnable example for CompileBinInst.
func ExampleCompileBinInst() {
	fmt.Println(recipe.CompileBinInst())

	// Output:
	// CompileBin
}

// ExamplePushInputInst is a runnable example for PushInputInst.
func ExamplePushInputInst() {
	fmt.Println(recipe.PushInputInst("foo.c"))

	// Output:
	// PushInput "foo.c"
}

// ExamplePushInputsInst is a runnable example for PushInputInst.
func ExamplePushInputsInst() {
	fmt.Println(recipe.PushInputsInst(filekind.C))

	// Output:
	// PushInputs "c"
}
