// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package transfer_test

import (
	"path"
	"reflect"
	"testing"

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
		maps map[string]string
	}{
		"empty": {
			in:   subject.Subject{},
			out:  subject.Subject{},
			maps: map[string]string{},
		},
		"litmus": {
			in:  subject.Subject{Litmus: path.Join("foo", "bar", "baz.litmus")},
			out: subject.Subject{Litmus: transfer.FileOrigLitmus},
			maps: map[string]string{
				transfer.FileOrigLitmus: path.Join("foo", "bar", "baz.litmus"),
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
			maps: map[string]string{
				transfer.FileFuzzLitmus: path.Join("barbaz", "baz.1.litmus"),
				transfer.FileFuzzTrace:  path.Join("barbaz", "baz.1.trace"),
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
			maps: map[string]string{
				path.Join(transfer.DirHarnesses, "arm", "inky.c"):  path.Join("burble", "armv8", "inky.c"),
				path.Join(transfer.DirHarnesses, "arm", "pinky.c"): path.Join("burble", "armv8", "pinky.c"),
				path.Join(transfer.DirHarnesses, "x86", "inky.c"):  path.Join("burble", "i386", "inky.c"),
				path.Join(transfer.DirHarnesses, "x86", "pinky.c"): path.Join("burble", "i386", "pinky.c"),
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
			maps: map[string]string{
				path.Join(transfer.DirCompiles, "clang", transfer.FileBin):        path.Join("foobaz", "clang", "a.out"),
				path.Join(transfer.DirCompiles, "gcc", transfer.FileBin):          path.Join("foobaz", "gcc", "a.out"),
				path.Join(transfer.DirCompiles, "clang", transfer.FileCompileLog): path.Join("foobaz", "clang", "errors"),
				path.Join(transfer.DirCompiles, "gcc", transfer.FileCompileLog):   path.Join("foobaz", "gcc", "errors"),
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			n := transfer.NewNormaliser("")
			s, err := n.Subject(c.in)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if !reflect.DeepEqual(*s, c.out) {
				t.Errorf("unexpected subject: got=%v, want=%v", s, c.out)
			}
			if !reflect.DeepEqual(n.Mappings, c.maps) {
				t.Errorf("unexpected mappings: got=%v, want=%v", n.Mappings, c.maps)
			}
		})
	}
}
