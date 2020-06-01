// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compile

// Kind is the enumeration of kinds of single-file compilation output.
type Kind uint8

const (
	// Exe refers to executable binary compilations.
	Exe Kind = iota
	// Obj refers to object file compilations.
	Obj
)

//go:generate stringer -type Kind
