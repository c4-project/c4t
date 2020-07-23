// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"io"

	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Factory is the interface of factories for machine node runners.
//
// Runner factories can contain disposable state (for example, long-running SSH connections), and so can be closed.
type Factory interface {
	// MakeRunner creates a new Runner, representing a particular invoker session on a machine.
	// It takes the plan in case the factory is waiting to get machine configuration from it.
	MakeRunner(p *plan.Plan, obs ...copier.Observer) (Runner, error)

	// Runner spawners can be closed once no more runners are needed.
	// For SSH runner spawners, this will close the SSH connection.
	io.Closer
}
