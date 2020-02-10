package ux

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

const usagePlanFile = "Read from this plan `file` instead of stdin."

// PlanLoaderFlags sets up a standard set of arguments feeding into the PlanLoader p.
func PlanLoaderFlags(p *model.PlanLoader) {
	flag.StringVar(&p.PlanFile, "i", "", usagePlanFile)
}
