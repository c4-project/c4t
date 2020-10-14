// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"context"
	"os"
	"strconv"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/model/job"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// Fuzz wraps the ACT one-file fuzzer, supplying the given seed.
func (a *Runner) Fuzz(ctx context.Context, j job.Fuzzer) error {
	seedStr := strconv.Itoa(int(j.Seed))

	cf, err := MakeFuzzConfFile(j)
	if err != nil {
		return err
	}

	args := []string{"-config", cf, "-seed", seedStr, "-o", j.OutLitmus}
	if ystring.IsNotEmpty(j.OutTrace) {
		args = append(args, "-trace-output", j.OutTrace, j.In)
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
