// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"io"

	c "github.com/urfave/cli/v2"
)

// SetCommonAppSettings sets app settings on a that are common to all act-tester endpoints.
// It passes through the pointer a.
func SetCommonAppSettings(a *c.App, outw, errw io.Writer) *c.App {
	a.Writer = outw
	a.ErrWriter = errw
	a.HideHelpCommand = true
	a.UseShortOptionHandling = true
	return a
}
