// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"encoding/csv"
	"io"

	"github.com/1set/gut/ystring"
	"github.com/c4-project/c4t/internal/helper/errhelp"
	"github.com/c4-project/c4t/internal/stat"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the analyser binary.
	Name  = "c4t-stat"
	usage = "inspects the statistics file"

	readme = `
   This program reads the statistics file maintained by the director, and
   prints CSV or human-readable summaries of its contents.`

	flagCsvMutations   = "csv-mutations"
	usageCsvMutations  = "dump CSV of mutation testing results"
	flagUseTotals      = "use-totals"
	flagUseTotalsShort = "t"
	usageUseTotals     = "use multi-session totals rather than per-session totals"
	flagStatFile       = "input"
	flagStatFileShort  = "i"
	usageStatFile      = "read statistics from this `FILE`"
)

// App is the entry point for c4t-analyse.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: readme,
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.ConfFileCliFlag(),
		&c.BoolFlag{Name: flagCsvMutations, Usage: usageCsvMutations},
		&c.BoolFlag{Name: flagUseTotals, Aliases: []string{flagUseTotalsShort}, Usage: usageUseTotals},
		&c.PathFlag{
			Name:        flagStatFile,
			Aliases:     []string{flagStatFileShort},
			Usage:       usageStatFile,
			DefaultText: "read from configuration",
		},
	}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	// TODO(@MattWindsor91): maybe use stat persister?
	f, err2 := getStatFile(ctx)
	if err2 != nil {
		return err2
	}
	var set stat.Set
	if err := set.Load(f); err != nil {
		return err
	}
	err := dump(ctx, &set, outw)
	cerr := f.Close()
	return errhelp.FirstError(err, cerr)
}

func dump(ctx *c.Context, set *stat.Set, w io.Writer) error {
	totals := ctx.Bool(flagUseTotals)
	csvMutations := ctx.Bool(flagCsvMutations)

	if csvMutations {
		if err := dumpCsvMutations(w, set, totals); err != nil {
			return err
		}
	}

	return nil
}

func dumpCsvMutations(w io.Writer, set *stat.Set, totals bool) error {
	cw := csv.NewWriter(w)
	if err := set.DumpMutationCSVHeader(cw); err != nil {
		return err
	}
	return set.DumpMutationCSV(cw, totals)
}

func getStatFile(ctx *c.Context) (io.ReadCloser, error) {
	fname, err := getStatPath(ctx)
	if err != nil {
		return nil, err
	}
	f, err := stat.OpenStatFile(fname)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// getStatPath computes the intended path to the stats file.
func getStatPath(ctx *c.Context) (string, error) {
	if f := ctx.Path(flagStatFile); ystring.IsNotBlank(f) {
		return f, nil
	}
	cfg, err := stdflag.ConfFileFromCli(ctx)
	if err != nil {
		return "", err
	}
	return cfg.Paths.StatFile()
}
