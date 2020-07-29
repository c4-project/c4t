// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"os"

	"github.com/MattWindsor91/act-tester/internal/app/invoke"

	"github.com/MattWindsor91/act-tester/internal/ux"
)

func main() {
	app := invoke.App(os.Stdout, os.Stderr)
	ux.LogTopError(app.Run(os.Args))
}
