// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"errors"
	"os"
)

// ErrPathsetNil is a standard error for when things that expect a pathset don't get one.
var ErrPathsetNil = errors.New("pathset nil")

// Mkdirs tries to make each directory in the directory set d.
func Mkdirs(ps ...string) error {
	for _, dir := range ps {
		if err := os.MkdirAll(dir, 0744); err != nil {
			return err
		}
	}
	return nil
}

// Rmdirs tries to remove each directory in the directory set d.
func Rmdirs(ps ...string) error {
	for _, dir := range ps {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}
