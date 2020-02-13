package main

import (
	"context"
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	// direct is the Director being built and run by this command.
	var (
		direct director.Director
		pf     string
	)

	ux.PlanFileFlag(&pf)
	flag.Parse()

	err := direct.Direct(context.Background())
	ux.LogTopError(err)
}
