// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

import (
	"io"
	"io/ioutil"

	"github.com/MattWindsor91/act-tester/internal/stage/analyse/pretty"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/stage/analyse/saver"

	"github.com/MattWindsor91/act-tester/internal/stage/analyse"

	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	name  = "act-tester-analyse"
	usage = "analyses a plan file"

	readme = `
   This program performs analysis on a plan file, and acts upon it.
   Analysis includes, at time of writing:

   - computing basic statistics on compile and run times per compiler;
   - categorising subjects by their final status.

   The program can act on its analysis in various ways, depending on the given
   flags.  By passing one or more -show flags, one can receive a human-readable
   summary of the plan file.  By passing -` + flagSaveDir + `, one can
   archive failing corpora to a directory for later experimentation.`

	flagShowCompilers      = "show-compilers"
	flagShowCompilersShort = "C"
	usageShowCompilers     = "show breakdown of compilers and their run times"
	flagShowOk             = "show-ok"
	flagShowOkShort        = "O"
	usageShowOk            = "show subjects that did not have compile or run issues"
	flagShowPlanInfo       = "show-plan-info"
	flagShowPlanInfoShort  = "P"
	usageShowPlanInfo      = "show plan metadata and stage times"
	flagShowSubjects       = "show-subjects"
	flagShowSubjectsShort  = "S"
	usageShowSubjects      = "show subjects by status"
	flagSaveDir            = "save-dir"
	usageSaveDir           = "if present, save failing corpora to this `directory`"
)

func App(outw, errw io.Writer) *c.App {
	a := &c.App{
		Name:        name,
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
		stdflag.WorkerCountCliFlag(),
		&c.BoolFlag{Name: flagShowCompilers, Aliases: []string{flagShowCompilersShort}, Usage: usageShowCompilers},
		&c.BoolFlag{Name: flagShowOk, Aliases: []string{flagShowOkShort}, Usage: usageShowOk},
		&c.BoolFlag{Name: flagShowSubjects, Aliases: []string{flagShowSubjectsShort}, Usage: usageShowSubjects},
		&c.BoolFlag{Name: flagShowPlanInfo, Aliases: []string{flagShowPlanInfoShort}, Usage: usageShowPlanInfo},
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
	obs, err := pretty.NewPrinter(
		pretty.WriteTo(outw),
		pretty.ShowCompilers(ctx.Bool(flagShowCompilers)),
		pretty.ShowOk(ctx.Bool(flagShowOk)),
		pretty.ShowSubjects(ctx.Bool(flagShowSubjects)),
		pretty.ShowPlanInfo(ctx.Bool(flagShowPlanInfo)),
	)
	if err != nil {
		return err
	}

	a, err := analyse.New(
		analyse.ObserveWith(obs),
		analyse.ParWorkers(stdflag.WorkerCountFromCli(ctx)),
		analyse.SaveToPathset(savedPaths(ctx)),
	)
	if err != nil {
		return err
	}
	return view.RunOnCliPlan(ctx, a, ioutil.Discard)
}

func savedPaths(ctx *c.Context) *saver.Pathset {
	root := ctx.Path(flagSaveDir)
	if ystring.IsBlank(root) {
		return nil
	}
	return saver.NewPathset(root)
}
