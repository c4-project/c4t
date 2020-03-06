// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"context"
	"flag"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

func main() {
	err := run(os.Args)
	ux.LogTopError(err)
}

const usageConfFile = "The `file` from which to load the tester configuration."

func run(args []string) error {
	cfile := flag.String("C", "", usageConfFile)

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	c, err := config.Load(*cfile)
	if err != nil {
		return err
	}

	d, err := director.New(c)
	if err != nil {
		return err
	}
	return d.Direct(context.Background())
}
