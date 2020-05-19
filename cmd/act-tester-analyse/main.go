// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"os"

	"github.com/MattWindsor91/act-tester/internal/app/analyse"
	"github.com/MattWindsor91/act-tester/internal/view"
)

func main() {
	view.LogTopError(analyse.App(os.Stdout, os.Stderr).Run(os.Args))
}
