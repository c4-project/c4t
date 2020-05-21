// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"
	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/MattWindsor91/act-tester/internal/model/status"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// ExampleMap_RenamesMatching is a runnable example for RenamesMatching.
func ExampleMap_RenamesMatching() {
	n := normaliser.New("root")
	s := subject.Subject{
		OrigLitmus: path.Join("foo", "bar", "baz.litmus"),
		Fuzz: &subject.Fuzz{
			Files: subject.FuzzFileset{
				Litmus: path.Join("barbaz", "baz.1.litmus"),
				Trace:  path.Join("barbaz", "baz.1.trace"),
			},
		},
		Compiles: map[string]subject.CompileResult{
			"clang": {
				Result: subject.Result{Status: status.Ok},
				Files: subject.CompileFileset{
					Bin: path.Join("foobaz", "clang", "a.out"),
					Log: path.Join("foobaz", "clang", "errors"),
				},
			},
		},
		Harnesses: map[string]subject.Harness{
			"arm": {
				Dir:   path.Join("burble", "armv8"),
				Files: []string{"inky.c", "pinky.c"},
			},
			"x86": {
				Dir:   path.Join("burble", "i386"),
				Files: []string{"inky.c", "pinky.c"},
			},
		},
	}
	_, _ = n.Normalise(s)
	for k, v := range n.Mappings.RenamesMatching(filekind.Any, filekind.InHarness) {
		fmt.Println(k, "<-", v)
	}

	// Unordered output:
	// root/harnesses/arm/inky.c <- burble/armv8/inky.c
	// root/harnesses/arm/pinky.c <- burble/armv8/pinky.c
	// root/harnesses/x86/inky.c <- burble/i386/inky.c
	// root/harnesses/x86/pinky.c <- burble/i386/pinky.c
}
