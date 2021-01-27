// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/recipe"
)

// ExampleRecipe_Paths is a testable example for Recipe.Paths.
func ExampleRecipe_Paths() {
	h, _ := recipe.New("foo/bar", recipe.OutNothing, recipe.AddFiles("baz", "barbaz", "foobar"))
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}

// ExampleRecipe_NeedsCompile is a testable example for Recipe.NeedsCompile.
func ExampleRecipe_NeedsCompile() {
	r1, _ := recipe.New("foo/bar", recipe.OutNothing, recipe.AddFiles("baz", "barbaz", "foobar"))
	fmt.Println("first recipe needs compile:", r1.NeedsCompile())
	r2, _ := recipe.New(
		"foo/bar",
		recipe.OutExe,
		recipe.AddFiles("baz", "barbaz", "foobar"),
		recipe.CompileAllCToExe(),
	)
	fmt.Println("second recipe needs compile:", r2.NeedsCompile())

	// Output:
	// first recipe needs compile: false
	// second recipe needs compile: true
}

// TestOp_UnmarshalJSON_error tests error cases of Op.UnmarshalJSON.
func TestOp_UnmarshalJSON_error(t *testing.T) {
	t.Parallel()

	var op recipe.Op

	dec := json.NewDecoder(strings.NewReader("6"))
	require.Error(t, dec.Decode(&op), "should not be able to decode this")
}
