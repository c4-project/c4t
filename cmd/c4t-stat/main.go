// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"os"

	"github.com/c4-project/c4t/internal/app/stat"

	"github.com/c4-project/c4t/internal/ux"
)

func main() {
	ux.LogTopError(stat.App(os.Stdout, os.Stderr).Run(os.Args))
}
