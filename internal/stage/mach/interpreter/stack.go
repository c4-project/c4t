// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter

import "errors"

// ErrUnderflow occurs if the interpreter's file stack underflows.
var ErrUnderflow = errors.New("stack underflow")

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

func (s *stack) popOne() (string, error) {
	pops := s.pop(1)
	if len(pops) != 1 {
		return "", ErrUnderflow
	}
	return pops[0], nil
}
