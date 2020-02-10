package main

import (
	"context"
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const usagePlanFile = "Read from this plan `file` instead of stdin."

// direct is the Director being built and run by this command.
var direct director.Director

func init() {
	ux.PlanLoaderFlags(&direct.PlanLoader)
}

func main() {
	flag.Parse()
	err := direct.Direct(context.Background())
	ux.LogTopError(err)
}
