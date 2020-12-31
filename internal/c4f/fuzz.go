// Copyright (c) 2020 Matt Windsor and contributors
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

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// Fuzz wraps the ACT one-file fuzzer, supplying the given seed.
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
		Cmd:    BinActFuzz,
		Subcmd: "run",
		Args:   args,
	}
	cerr := a.Run(ctx, cs)
	rerr := os.Remove(cf)
	return errhelp.FirstError(cerr, rerr)
}
