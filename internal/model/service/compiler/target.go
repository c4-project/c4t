// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

// Target is the enumeration of kinds of single-file compilation output.
type Target uint8

const (
	// Exe refers to executable binary compilations.
	Exe Target = iota
	// Obj refers to object file compilations.
	Obj
)

//go:generate stringer -type Target
