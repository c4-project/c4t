package main

import (
	"context"
	"flag"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	err := run(os.Args)
	ux.LogTopError(err)
}

func run(args []string) error {
	var (
		// direct is the Director being built and run by this command.
		direct director.Director
		pf     string
	)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	ux.PlanFileFlag(fs, &pf)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return direct.Direct(context.Background())
}
