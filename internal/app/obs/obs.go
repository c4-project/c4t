// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/c4t/internal/helper/errhelp"
	obs2 "github.com/MattWindsor91/c4t/internal/subject/obs"
	"github.com/MattWindsor91/c4t/internal/ux"

	"github.com/MattWindsor91/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the backend binary.
	Name  = "c4t-obs"
	usage = "interpret observation files"

	readme = `
    This program interprets observation JSON files of the form produced by act-backend.`

	flagShowPostcondition      = "show-postcondition"
	flagShowPostconditionShort = "p"
	usageShowPostcondition     = "print a Litmus (forall, sum of products) postcondition capturing the states observed"
)

// App is the c4-obs app.
func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        Name,
		Usage:       usage,
		Description: readme,
		Flags:       flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		&c.BoolFlag{
			Name:    flagShowPostcondition,
			Aliases: []string{flagShowPostconditionShort},
			Usage:   usageShowPostcondition,
		},
	}
}

func run(ctx *c.Context, outw io.Writer) error {
	m := modeFromCli(ctx)

	var o obs2.Obs
	if err := readObs(ctx, &o); err != nil {
		return err
	}

	return obs2.Pretty(outw, o, m)
}

func modeFromCli(ctx *c.Context) obs2.PrettyMode {
	return obs2.PrettyMode{
		Dnf: ctx.Bool(flagShowPostcondition),
	}
}

func readObs(ctx *c.Context, o *obs2.Obs) error {
	fname, err := fileFromCli(ctx)
	if err != nil {
		return err
	}
	r, err := openFileOrStdin(fname)
	if err != nil {
		return err
	}
	derr := json.NewDecoder(r).Decode(o)
	cerr := r.Close()
	return errhelp.FirstError(derr, cerr)
}

func fileFromCli(ctx *c.Context) (string, error) {
	// TODO(@MattWindsor91): generalise this.
	nargs := ctx.Args().Len()
	switch nargs {
	case 0:
		return "", nil
	case 1:
		return ctx.Args().First(), nil
	default:
		return "", fmt.Errorf("too many arguments: expected 1, got %d", nargs)
	}
}

func openFileOrStdin(fname string) (io.ReadCloser, error) {
	if ystring.IsBlank(fname) || fname == ux.StdinFile {
		return os.Stdin, nil
	}
	return os.Open(fname)
}
