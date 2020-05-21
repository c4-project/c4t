// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

import (
	"io"
	"io/ioutil"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse"

	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	flagSaveDir  = "save-dir"
	usageSaveDir = "if present, save failing corpora to this `directory`"
)

func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:  "act-tester-analyse",
		Usage: "performs human-readable queries on a plan file",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
	}
	return stdflag.SetPlanAppSettings(a, outw, errw)
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.WorkerCountCliFlag(),
		&c.PathFlag{
			Name:        flagSaveDir,
			Aliases:     []string{stdflag.FlagOutDir},
			Usage:       usageSaveDir,
			DefaultText: "do not save",
		},
		// TODO(@MattWindsor91): template stuff
	}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	obs, err := observer.NewAnalysisWriter(outw)
	if err != nil {
		return err
	}

	q := analyse.Config{Observers: []observer.Observer{obs},
		NWorkers:   stdflag.WorkerCountFromCli(ctx),
		SavedPaths: savedPaths(ctx),
	}
	return view.RunOnCliPlan(ctx, &q, ioutil.Discard)
}

func savedPaths(ctx *c.Context) *saver.Pathset {
	root := ctx.Path(flagSaveDir)
	if ystring.IsBlank(root) {
		return nil
	}
	return saver.NewPathset(root)
}
