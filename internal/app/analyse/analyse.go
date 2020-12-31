// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

import (
	"io"
	"io/ioutil"

	"github.com/c4-project/c4t/internal/plan/analysis"

	"github.com/c4-project/c4t/internal/stage/analyser/csvdump"
	"github.com/c4-project/c4t/internal/stage/analyser/pretty"

	"github.com/1set/gut/ystring"
	"github.com/c4-project/c4t/internal/stage/analyser/saver"

	"github.com/c4-project/c4t/internal/stage/analyser"

	"github.com/c4-project/c4t/internal/ux"

	"github.com/c4-project/c4t/internal/ux/stdflag"
	c "github.com/urfave/cli/v2"
)

const (
	// Name is the name of the analyser binary.
	Name  = "c4t-analyse"
	usage = "analyses a plan file"

	readme = `
   This program performs analysis on a plan file, and acts upon it.
   Analysis includes, at time of writing:

   - computing basic statistics on compile and run times per compiler;
   - categorising subjects by their final status.

   The program can c4f on its analysis in various ways, depending on the given
   flags.  By passing one or more -show flags, one can receive a human-readable
   summary of the plan file.  By passing -` + flagSaveDir + `, one can
   archive failing corpora to a directory for later experimentation.`

	// FlagErrorOnBadStatus is used to activate error-on-bad-status.  It is exported for testing purposes.
	FlagErrorOnBadStatus      = "error-on-bad-status"
	flagErrorOnBadStatusShort = "e"
	usageErrorOnBadStatus     = "report an error if plan contains subjects with bad statuses"
	flagLoadFilters           = "filter-file"
	usageLoadFilters          = "load compile result filters from this file"
	flagCsvCompilers          = "csv-compilers"
	usageCsvCompilers         = "dump CSV of compilers and their run times"
	flagCsvStages             = "csv-stages"
	usageCsvStages            = "dump CSV of stages and their run times"
	flagShowCompilers         = "show-compilers"
	flagShowCompilersShort    = "C"
	usageShowCompilers        = "show breakdown of compilers and their run times"
	flagShowCompilerLogs      = "show-compiler-logs"
	flagShowCompilerLogsShort = "L"
	usageShowCompilerLogs     = "show breakdown of compiler logs (requires -" + flagShowCompilers + ")"
	flagShowOk                = "show-ok"
	flagShowOkShort           = "O"
	usageShowOk               = "show subjects that did not have compile or run issues"
	flagShowPlanInfo          = "show-plan-info"
	flagShowPlanInfoShort     = "P"
	usageShowPlanInfo         = "show plan metadata and stage times"
	flagShowSubjects          = "show-subjects"
	flagShowSubjectsShort     = "S"
	usageShowSubjects         = "show subjects by status"
	flagSaveDir               = "save-dir"
	usageSaveDir              = "if present, save failing corpora to this `directory`"
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
		stdflag.WorkerCountCliFlag(),
		&c.BoolFlag{Name: FlagErrorOnBadStatus, Aliases: []string{flagErrorOnBadStatusShort}, Usage: usageErrorOnBadStatus},
		&c.BoolFlag{Name: flagCsvCompilers, Usage: usageCsvCompilers},
		&c.BoolFlag{Name: flagCsvStages, Usage: usageCsvStages},
		&c.BoolFlag{Name: flagShowCompilers, Aliases: []string{flagShowCompilersShort}, Usage: usageShowCompilers},
		&c.BoolFlag{Name: flagShowCompilerLogs, Aliases: []string{flagShowCompilerLogsShort}, Usage: usageShowCompilerLogs},
		&c.BoolFlag{Name: flagShowOk, Aliases: []string{flagShowOkShort}, Usage: usageShowOk},
		&c.BoolFlag{Name: flagShowSubjects, Aliases: []string{flagShowSubjectsShort}, Usage: usageShowSubjects},
		&c.BoolFlag{Name: flagShowPlanInfo, Aliases: []string{flagShowPlanInfoShort}, Usage: usageShowPlanInfo},
		&c.PathFlag{
			Name:        flagSaveDir,
			Aliases:     []string{stdflag.FlagOutDir},
			Usage:       usageSaveDir,
			DefaultText: "do not save",
		},
		&c.PathFlag{
			Name:        flagLoadFilters,
			Usage:       usageLoadFilters,
			DefaultText: "do not load filters",
		},
	}
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	obs, err := observers(ctx, outw)
	if err != nil {
		return err
	}

	a, err := analyser.New(
		analyser.ObserveWith(obs...),
		analyser.Analysis(
			analysis.WithFiltersFromFile(ctx.Path(flagLoadFilters)),
			analysis.WithWorkerCount(stdflag.WorkerCountFromCli(ctx)),
		),
		analyser.ErrorOnBadStatus(ctx.Bool(FlagErrorOnBadStatus)),
		analyser.SaveToPathset(savedPaths(ctx)),
	)
	if err != nil {
		return err
	}
	return ux.RunOnCliPlan(ctx, a, ioutil.Discard)
}

func observers(ctx *c.Context, outw io.Writer) ([]analyser.Observer, error) {
	obs, err := prettyObserver(ctx, outw)
	if err != nil {
		return nil, err
	}
	return csvObserver(ctx, outw, obs)
}

func prettyObserver(ctx *c.Context, outw io.Writer) ([]analyser.Observer, error) {
	showCompilers := ctx.Bool(flagShowCompilers)
	// showCompilerLogs depends on showCompilers
	showOk := ctx.Bool(flagShowOk)
	showSubjects := ctx.Bool(flagShowSubjects)
	showPlanInfo := ctx.Bool(flagShowPlanInfo)

	if !(showCompilers || showOk || showSubjects || showPlanInfo) {
		return nil, nil
	}
	po, err := pretty.NewPrinter(
		pretty.WriteTo(outw),
		pretty.ShowCompilers(showCompilers),
		pretty.ShowCompilerLogs(ctx.Bool(flagShowCompilerLogs)),
		pretty.ShowOk(showOk),
		pretty.ShowSubjects(showSubjects),
		pretty.ShowPlanInfo(showPlanInfo),
	)
	return []analyser.Observer{po}, err
}

func csvObserver(ctx *c.Context, outw io.Writer, obs []analyser.Observer) ([]analyser.Observer, error) {
	if ctx.Bool(flagCsvCompilers) {
		obs = append(obs, csvdump.NewCompilerWriter(outw))
	}
	if ctx.Bool(flagCsvStages) {
		obs = append(obs, csvdump.NewStageWriter(outw))
	}
	return obs, nil
}

func savedPaths(ctx *c.Context) *saver.Pathset {
	root := ctx.Path(flagSaveDir)
	if ystring.IsBlank(root) {
		return nil
	}
	return saver.NewPathset(root)
}
