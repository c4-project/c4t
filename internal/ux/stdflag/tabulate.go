// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

import (
	"io"

	"github.com/c4-project/c4t/internal/tabulator"
	c "github.com/urfave/cli/v2"
)

const (
	// FlagOutputCsv is a flag for selecting the CSV tabulator.
	FlagOutputCsv = "output-csv"
	// UsageOutputCsv is the usage for the output-CSV flag.
	UsageOutputCsv = "output tables as CSV"
)

// TabulatorCliFlags expands into the flag set for configuring TabulatorFromCli.
func TabulatorCliFlags() []c.Flag {
	return []c.Flag{
		&c.BoolFlag{
			Name:  FlagOutputCsv,
			Usage: UsageOutputCsv,
		},
	}
}

// TabulatorFromCli constructs a Tabulator that corresponds to the flags set in ctx.
//
// There is only one format choice at the moment (CSV vs tab), but this may change in future.
func TabulatorFromCli(ctx *c.Context, w io.Writer) tabulator.Tabulator {
	if ctx.Bool(FlagOutputCsv) {
		return tabulator.NewCsv(w)
	}
	return tabulator.NewTab(w)
}
