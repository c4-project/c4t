// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package normalise_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/normalise"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// TODO(@MattWindsor91): test rooting

// TestNormaliser_Subject checks the normaliser on various small subject cases.
func TestNormaliser_Subject(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   subject.Subject
		out  subject.Subject
		maps map[string]normalise.Normalisation
	}{
		"empty": {
			in:   subject.Subject{},
			out:  subject.Subject{},
			maps: map[string]normalise.Normalisation{},
		},
		"litmus": {
			in:  subject.Subject{OrigLitmus: path.Join("foo", "bar", "baz.litmus")},
			out: subject.Subject{OrigLitmus: normalise.FileOrigLitmus},
			maps: map[string]normalise.Normalisation{
				normalise.FileOrigLitmus: {
					Original: path.Join("foo", "bar", "baz.litmus"),
					Kind:     filekind.Litmus,
					Loc:      filekind.InOrig,
				},
			},
		},
		"fuzz": {
			in: subject.Subject{
				Fuzz: &subject.Fuzz{
					Files: subject.FuzzFileset{
						Litmus: path.Join("barbaz", "baz.1.litmus"),
						Trace:  path.Join("barbaz", "baz.1.trace"),
					},
				},
			},
			out: subject.Subject{
				Fuzz: &subject.Fuzz{
					Files: subject.FuzzFileset{
						Litmus: normalise.FileFuzzLitmus,
						Trace:  normalise.FileFuzzTrace,
					},
				},
			},
			maps: map[string]normalise.Normalisation{
				normalise.FileFuzzLitmus: {
					Original: path.Join("barbaz", "baz.1.litmus"),
					Kind:     filekind.Litmus,
					Loc:      filekind.InFuzz,
				},
				normalise.FileFuzzTrace: {
					Original: path.Join("barbaz", "baz.1.trace"),
					Kind:     filekind.Trace,
					Loc:      filekind.InFuzz,
				},
			},
		},
		"harness": {
			in: subject.Subject{
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
			},
			out: subject.Subject{
				Harnesses: map[string]subject.Harness{
					"arm": {
						Dir:   path.Join(normalise.DirHarnesses, "arm"),
						Files: []string{"inky.c", "pinky.c"},
					},
					"x86": {
						Dir:   path.Join(normalise.DirHarnesses, "x86"),
						Files: []string{"inky.c", "pinky.c"},
					},
				},
			},
			maps: map[string]normalise.Normalisation{
				path.Join(normalise.DirHarnesses, "arm", "inky.c"): {
					Original: path.Join("burble", "armv8", "inky.c"),
					Kind:     filekind.CSrc,
					Loc:      filekind.InHarness,
				},
				path.Join(normalise.DirHarnesses, "arm", "pinky.c"): {
					Original: path.Join("burble", "armv8", "pinky.c"),
					Kind:     filekind.CSrc,
					Loc:      filekind.InHarness,
				},
				path.Join(normalise.DirHarnesses, "x86", "inky.c"): {
					Original: path.Join("burble", "i386", "inky.c"),
					Kind:     filekind.CSrc,
					Loc:      filekind.InHarness,
				},
				path.Join(normalise.DirHarnesses, "x86", "pinky.c"): {
					Original: path.Join("burble", "i386", "pinky.c"),
					Kind:     filekind.CSrc,
					Loc:      filekind.InHarness,
				},
			},
		},
		"compile": {
			in: subject.Subject{
				Compiles: map[string]subject.CompileResult{
					"clang": {
						Result: subject.Result{Status: status.Ok},
						Files: subject.CompileFileset{
							Bin: path.Join("foobaz", "clang", "a.out"),
							Log: path.Join("foobaz", "clang", "errors"),
						},
					},
					"gcc": {
						Result: subject.Result{Status: status.Ok},
						Files: subject.CompileFileset{
							Bin: path.Join("foobaz", "gcc", "a.out"),
							Log: path.Join("foobaz", "gcc", "errors"),
						},
					},
				},
			},
			out: subject.Subject{
				Compiles: map[string]subject.CompileResult{
					"clang": {
						Result: subject.Result{Status: status.Ok},
						Files: subject.CompileFileset{
							Bin: path.Join(normalise.DirCompiles, "clang", normalise.FileBin),
							Log: path.Join(normalise.DirCompiles, "clang", normalise.FileCompileLog),
						},
					},
					"gcc": {
						Result: subject.Result{Status: status.Ok},
						Files: subject.CompileFileset{
							Bin: path.Join(normalise.DirCompiles, "gcc", normalise.FileBin),
							Log: path.Join(normalise.DirCompiles, "gcc", normalise.FileCompileLog),
						},
					},
				},
			},
			maps: map[string]normalise.Normalisation{
				path.Join(normalise.DirCompiles, "clang", normalise.FileBin): {
					Original: path.Join("foobaz", "clang", "a.out"),
					Kind:     filekind.Bin,
					Loc:      filekind.InCompile,
				},
				path.Join(normalise.DirCompiles, "gcc", normalise.FileBin): {
					Original: path.Join("foobaz", "gcc", "a.out"),
					Kind:     filekind.Bin,
					Loc:      filekind.InCompile,
				},
				path.Join(normalise.DirCompiles, "clang", normalise.FileCompileLog): {
					Original: path.Join("foobaz", "clang", "errors"),
					Kind:     filekind.Log,
					Loc:      filekind.InCompile,
				},
				path.Join(normalise.DirCompiles, "gcc", normalise.FileCompileLog): {
					Original: path.Join("foobaz", "gcc", "errors"),
					Kind:     filekind.Log,
					Loc:      filekind.InCompile,
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			n := normalise.NewNormaliser("")
			s, err := n.Subject(c.in)
			if assert.NoError(t, err) {
				assert.Equal(t, c.out, *s)
				assert.Equal(t, c.maps, n.Mappings)
			}
		})
	}
}

// ExampleNormaliser_MappingsMatching is a runnable example for MappingsMatching.
func ExampleNormaliser_MappingsMatching() {
	n := normalise.NewNormaliser("root")
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
	_, _ = n.Subject(s)
	for k, v := range n.MappingsMatching(filekind.Any, filekind.InHarness) {
		fmt.Println(k, "<-", v)
	}

	// Unordered output:
	// root/harnesses/arm/inky.c <- burble/armv8/inky.c
	// root/harnesses/arm/pinky.c <- burble/armv8/pinky.c
	// root/harnesses/x86/inky.c <- burble/i386/inky.c
	// root/harnesses/x86/pinky.c <- burble/i386/pinky.c
}
