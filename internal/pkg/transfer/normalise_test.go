// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package transfer_test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// TODO(@MattWindsor91): test rooting

// TestNormaliser_Subject checks the normaliser on various small subject cases.
func TestNormaliser_Subject(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   subject.Subject
		out  subject.Subject
		maps map[string]transfer.Normalisation
	}{
		"empty": {
			in:   subject.Subject{},
			out:  subject.Subject{},
			maps: map[string]transfer.Normalisation{},
		},
		"litmus": {
			in:  subject.Subject{Litmus: path.Join("foo", "bar", "baz.litmus")},
			out: subject.Subject{Litmus: transfer.FileOrigLitmus},
			maps: map[string]transfer.Normalisation{
				transfer.FileOrigLitmus: {
					Original: path.Join("foo", "bar", "baz.litmus"),
					Kind:     transfer.NKOrigLitmus,
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
						Litmus: transfer.FileFuzzLitmus,
						Trace:  transfer.FileFuzzTrace,
					},
				},
			},
			maps: map[string]transfer.Normalisation{
				transfer.FileFuzzLitmus: {
					Original: path.Join("barbaz", "baz.1.litmus"),
					Kind:     transfer.NKFuzz,
				},
				transfer.FileFuzzTrace: {
					Original: path.Join("barbaz", "baz.1.trace"),
					Kind:     transfer.NKFuzz,
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
						Dir:   path.Join(transfer.DirHarnesses, "arm"),
						Files: []string{"inky.c", "pinky.c"},
					},
					"x86": {
						Dir:   path.Join(transfer.DirHarnesses, "x86"),
						Files: []string{"inky.c", "pinky.c"},
					},
				},
			},
			maps: map[string]transfer.Normalisation{
				path.Join(transfer.DirHarnesses, "arm", "inky.c"): {
					Original: path.Join("burble", "armv8", "inky.c"),
					Kind:     transfer.NKHarness,
				},
				path.Join(transfer.DirHarnesses, "arm", "pinky.c"): {
					Original: path.Join("burble", "armv8", "pinky.c"),
					Kind:     transfer.NKHarness,
				},
				path.Join(transfer.DirHarnesses, "x86", "inky.c"): {
					Original: path.Join("burble", "i386", "inky.c"),
					Kind:     transfer.NKHarness,
				},
				path.Join(transfer.DirHarnesses, "x86", "pinky.c"): {
					Original: path.Join("burble", "i386", "pinky.c"),
					Kind:     transfer.NKHarness,
				},
			},
		},
		"compile": {
			in: subject.Subject{
				Compiles: map[string]subject.CompileResult{
					"clang": {
						Success: true,
						Files: subject.CompileFileset{
							Bin: path.Join("foobaz", "clang", "a.out"),
							Log: path.Join("foobaz", "clang", "errors"),
						},
					},
					"gcc": {
						Success: true,
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
						Success: true,
						Files: subject.CompileFileset{
							Bin: path.Join(transfer.DirCompiles, "clang", transfer.FileBin),
							Log: path.Join(transfer.DirCompiles, "clang", transfer.FileCompileLog),
						},
					},
					"gcc": {
						Success: true,
						Files: subject.CompileFileset{
							Bin: path.Join(transfer.DirCompiles, "gcc", transfer.FileBin),
							Log: path.Join(transfer.DirCompiles, "gcc", transfer.FileCompileLog),
						},
					},
				},
			},
			maps: map[string]transfer.Normalisation{
				path.Join(transfer.DirCompiles, "clang", transfer.FileBin): {
					Original: path.Join("foobaz", "clang", "a.out"),
					Kind:     transfer.NKCompile,
				},
				path.Join(transfer.DirCompiles, "gcc", transfer.FileBin): {
					Original: path.Join("foobaz", "gcc", "a.out"),
					Kind:     transfer.NKCompile,
				},
				path.Join(transfer.DirCompiles, "clang", transfer.FileCompileLog): {
					Original: path.Join("foobaz", "clang", "errors"),
					Kind:     transfer.NKCompile,
				},
				path.Join(transfer.DirCompiles, "gcc", transfer.FileCompileLog): {
					Original: path.Join("foobaz", "gcc", "errors"),
					Kind:     transfer.NKCompile,
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			n := transfer.NewNormaliser("")
			s, err := n.Subject(c.in)
			if assert.NoError(t, err) {
				assert.Equal(t, c.out, *s)
				assert.Equal(t, c.maps, n.Mappings)
			}
		})
	}
}
