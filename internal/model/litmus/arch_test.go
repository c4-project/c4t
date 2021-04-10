// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/litmus"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/id"
)

// ExampleArchToLitmus gives a few testable examples of ArchToLitmus.
func ExampleArchToLitmus() {
	a1, _ := litmus.ArchToLitmus(id.ArchAArch64)
	fmt.Println(a1)

	a2, _ := litmus.ArchToLitmus(id.ArchPPCPOWER9)
	fmt.Println(a2)

	a3, _ := litmus.ArchToLitmus(id.ArchX8664)
	fmt.Println(a3)

	a4, _ := litmus.ArchToLitmus(id.ArchC)
	fmt.Println(a4)

	// Output:
	// AArch64
	// PPC
	// X86_64
	// C
}

// ExampleArchOfLitmus gives a few testable examples of ArchOfLitmus.
func ExampleArchOfLitmus() {
	a1, _ := litmus.ArchOfLitmus("AArch64")
	fmt.Println(a1)

	a2, _ := litmus.ArchOfLitmus("PPC")
	fmt.Println(a2)

	a3, _ := litmus.ArchOfLitmus("X86_64")
	fmt.Println(a3)

	a4, _ := litmus.ArchOfLitmus("C")
	fmt.Println(a4)

	// Output:
	// aarch64
	// ppc
	// x86.64
	// c
}

// TestArchToLitmus_errors tests various failing cases of ArchOfLitmus.
func TestArchToLitmus_errors(t *testing.T) {
	t.Parallel()

	for name, c := range map[string]struct {
		in  id.ID
		err error
	}{
		"empty":   {err: litmus.ErrEmptyArch},
		"unknown": {in: id.FromString("notarealid"), err: litmus.ErrBadArch},
	} {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := litmus.ArchToLitmus(c.in)
			testhelp.ExpectErrorIs(t, err, c.err, "ArchOfLitmus")
		})
	}
}

// TestArchOfFile tests various happy-path attempts at getting architectures from Litmus files.
func TestArchOfFile(t *testing.T) {
	t.Parallel()

	for name, want := range map[string]id.ID{
		"aarch64": id.ArchAArch64,
		"c":       id.ArchC,
	} {
		name, want := name, want
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := litmus.ArchOfFile(filepath.Join("testdata", name+".litmus"))
			require.NoError(t, err, "arch deduction should work")
			require.Equal(t, want, got, "arch ID not as expected")
		})
	}
}

// TestArchOfFile_errors tests various failed attempts at getting architectures from Litmus files.
func TestArchOfFile_errors(t *testing.T) {
	t.Parallel()

	for name, want := range map[string]error{
		"blank":  litmus.ErrEmptyArch,
		"empty":  io.EOF,
		"bad":    litmus.ErrBadArch,
		"nofile": os.ErrNotExist,
	} {
		name, want := name, want
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := litmus.ArchOfFile(filepath.Join("testdata", name+".litmus"))
			testhelp.ExpectErrorIs(t, err, want, "running ArchOfFile")
		})
	}
}
