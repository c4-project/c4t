// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package recipe_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// ExampleRecipe_Paths is a testable example for Recipe.Paths.
func ExampleRecipe_Paths() {
	h := recipe.New("foo/bar", recipe.OutNothing, recipe.AddFiles("baz", "barbaz", "foobar"))
	for _, f := range h.Paths() {
		fmt.Println(f)
	}

	// Output:
	// foo/bar/baz
	// foo/bar/barbaz
	// foo/bar/foobar
}

func TestOp_UnmarshalJSON_error(t *testing.T) {
	t.Parallel()

	var op recipe.Op

	dec := json.NewDecoder(strings.NewReader("6"))
	require.Error(t, dec.Decode(&op), "should not be able to decode this")
}
