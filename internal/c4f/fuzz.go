// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f

import (
	"context"
	"os"
	"strconv"

	"github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

// BinC4fFuzz is the name of the c4f fuzzer binary.
const BinC4fFuzz = "c4f"

// Fuzz wraps the c4f one-file fuzzer, supplying the given seed.
func (a *Runner) Fuzz(ctx context.Context, j fuzzer.Job) error {
	cf, err := MakeFuzzConfFile(j)
	if err != nil {
		return err
	}

	args := []string{"-config", cf, "-seed", strconv.Itoa(int(j.Seed)), "-o", j.OutLitmus}
	if ystring.IsNotEmpty(j.OutTrace) {
		args = append(args, "-trace-output", j.OutTrace)
	}
	args = append(args, j.In)

	cs := CmdSpec{
		Cmd:    BinC4fFuzz,
		Subcmd: "run",
		Args:   args,
	}
	cerr := a.Run(ctx, cs)
	rerr := os.Remove(cf)
	return errhelp.FirstError(cerr, rerr)
}
