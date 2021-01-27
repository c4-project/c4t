// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package main

import (
	"os"

	"github.com/c4-project/c4t/internal/app/gccnt"

	"github.com/c4-project/c4t/internal/ux"
)

func main() {
	ux.LogTopError(gccnt.App(os.Stdout, os.Stderr).Run(os.Args))
}
