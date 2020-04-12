// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import "errors"

var (
	// ErrConfigNil occurs when the config passed into New is empty.
	ErrConfigNil = errors.New("config nil")

	// ErrDirEmpty occurs when the local directory filepath is empty.
	ErrDirEmpty = errors.New("local dir is empty string")

	// ErrSSHNil occurs when the SSH config part of a rmach config is nil.
	ErrSSHNil = errors.New("ssh config nil")
)
