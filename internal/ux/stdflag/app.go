// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"io"

	c "github.com/urfave/cli/v2"
)

const (
	at = "@"
	ic = "imperial.ac.uk"
)

var authors = [...]*c.Author{
	{Name: "Matt Windsor", Email: "m.windsor" + at + ic},
}

// SetCommonAppSettings sets app settings on a that are common to all act-tester endpoints.
// It passes through the pointer a.
func SetCommonAppSettings(a *c.App, outw, errw io.Writer) *c.App {
	a.Authors = authors[:]
	a.Writer = outw
	a.ErrWriter = errw
	a.HideHelpCommand = true
	a.UseShortOptionHandling = true
	return a
}

const planDescription = /* newline intentional */ `

   This command takes an optional argument naming the 'plan file' to load.  If
   this argument isn't present, or is set to '-', it will instead load the plan
   from stdin.`

// SetPlanAppSettings sets app settings on a that are common to act-tester endpoints handling a plan.
// It passes through the pointer a.
func SetPlanAppSettings(a *c.App, outw, errw io.Writer) *c.App {
	a.ArgsUsage = "[plan file]"
	a.Description += planDescription
	return SetCommonAppSettings(a, outw, errw)
}
