package ux

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

const (
	usageConfFile = "read ACT config from this `file`"
	usageDuneExec = "if true, use 'dune exec' to run OCaml ACT binaries"
	usagePlanFile = "read from this plan `file` instead of stdin"
)

// ActRunnerFlags sets up a standard set of arguments feeding into the ActRunner a.
func ActRunnerFlags(a *interop.ActRunner) {
	flag.StringVar(&a.ConfFile, "C", "", usageConfFile)
	flag.BoolVar(&a.DuneExec, "x", false, usageDuneExec)
}

// PlanLoaderFlags sets up a standard set of arguments feeding into the PlanLoader p.
func PlanLoaderFlags(p *model.PlanLoader) {
	flag.StringVar(&p.PlanFile, "i", "", usagePlanFile)
}
