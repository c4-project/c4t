// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter

// stack is the type of file-path stacks.
type stack []string

func (s *stack) push(file string) {
	*s = append(*s, file)
}

func (s *stack) pop(n int) []string {
	lfs := len(*s)
	if n <= 0 || lfs < n {
		n = lfs
	}
	cut := lfs - n

	var fs []string
	fs, *s = (*s)[cut:], (*s)[:cut]
	return fs
}
