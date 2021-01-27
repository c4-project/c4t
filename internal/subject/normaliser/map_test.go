// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"fmt"
	"path"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/model/filekind"
	"github.com/c4-project/c4t/internal/subject"
	"github.com/c4-project/c4t/internal/subject/normaliser"
	"github.com/c4-project/c4t/internal/subject/status"
)

// ExampleMap_RenamesMatching is a runnable example for RenamesMatching.
func ExampleMap_RenamesMatching() {
	n := normaliser.New("root")
	l, _ := litmus.New(path.Join("foo", "bar", "baz.litmus"))
	f, _ := litmus.New(path.Join("barbaz", "baz.1.litmus"))
	s, _ := subject.New(l,
		subject.WithFuzz(&subject.Fuzz{Litmus: *f, Trace: path.Join("barbaz", "baz.1.trace")}),
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
