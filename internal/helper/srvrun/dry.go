// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package srvrun

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// DryRunner is a runner that just outputs all commands to be run to itself.
type DryRunner struct {
	io.Writer
}

// Run dumps r to the dry runner's writer.
func (d DryRunner) Run(_ context.Context, r service.RunInfo) error {
	_, err := io.WriteString(d, r.String())
	return err
}
