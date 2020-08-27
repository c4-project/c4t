// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/subject/compilation"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
	"github.com/MattWindsor91/act-tester/internal/subject"
	"github.com/MattWindsor91/act-tester/internal/subject/normaliser"
	"github.com/MattWindsor91/act-tester/internal/subject/status"
)

// ExampleMap_RenamesMatching is a runnable example for RenamesMatching.
func ExampleMap_RenamesMatching() {
	n := normaliser.New("root")
	s, _ := subject.New(
		litmus.New(path.Join("foo", "bar", "baz.litmus")),
		subject.WithFuzz(
			&subject.Fuzz{
				Litmus: *litmus.New(path.Join("barbaz", "baz.1.litmus")),
				Trace:  path.Join("barbaz", "baz.1.trace"),
			},
		),
		subject.WithCompile(id.FromString("clang"),
			compilation.CompileResult{
				Result: compilation.Result{Status: status.Ok},
				Files: compilation.CompileFileset{
					Bin: path.Join("foobaz", "clang", "a.out"),
					Log: path.Join("foobaz", "clang", "errors"),
				},
			},
		),
		subject.WithRecipe(id.FromString("arm"),
			recipe.Recipe{
				Dir:   path.Join("burble", "armv8"),
				Files: []string{"inky.c", "pinky.c"},
			},
		),
		subject.WithRecipe(id.FromString("x86"),
			recipe.Recipe{
				Dir:   path.Join("burble", "i386"),
				Files: []string{"inky.c", "pinky.c"},
			},
		),
	)
	_, _ = n.Normalise(*s)
	for k, v := range n.Mappings.RenamesMatching(filekind.Any, filekind.InRecipe) {
		fmt.Println(k, "<-", v)
	}

	// Unordered output:
	// root/recipes/arm/inky.c <- burble/armv8/inky.c
	// root/recipes/arm/pinky.c <- burble/armv8/pinky.c
	// root/recipes/x86/inky.c <- burble/i386/inky.c
	// root/recipes/x86/pinky.c <- burble/i386/pinky.c
}
