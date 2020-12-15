// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"context"
	"os"

	"github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/helper/srvrun"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/litmus"

	"github.com/MattWindsor91/c4t/internal/model/id"
	mdl "github.com/MattWindsor91/c4t/internal/model/litmus"
	"github.com/MattWindsor91/c4t/internal/model/service"
)

// ExampleInstance_Run is a testable example for Run.
func ExampleInstance_Run() {
	i := litmus.Instance{
		Job: backend.LiftJob{
			Arch: id.ArchX8664,
			In:   backend.LiftLitmusInput(mdl.NewOrPanic("in.litmus", mdl.WithArch(id.ArchC))),
			Out:  backend.LiftOutput{Dir: "out", Target: backend.ToExeRecipe},
		},
		RunInfo: service.RunInfo{Cmd: "litmus7", Args: []string{"-v"}},
		Runner:  srvrun.DryRunner{Writer: os.Stdout},
	}

	// We don't ask for a fixset, so we won't have any patching.
	_ = i.Run(context.Background())

	// Output:
	// litmus7 -v -o out -carch X86_64 -c11 true in.litmus
}
