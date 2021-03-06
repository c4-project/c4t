// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"time"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/subject/corpus"
)

func MockMetadata() Metadata {
	return Metadata{
		Creation: time.Date(2011, time.November, 11, 11, 11, 11, 0, time.FixedZone("PST", -8*60*60)),
		Seed:     8675309,
		Version:  CurrentVer,
	}
}

// Mock makes a valid, but mocked-up, plan.
func Mock() *Plan {
	// TODO(@MattWindsor91): add things to this plan as time goes on.
	return &Plan{
		Metadata: MockMetadata(),
		Machine: machine.Named{
			ID: id.FromString("localhost"),
		},
		Backend: &backend2.NamedSpec{
			ID: id.FromString("litmus"),
			Spec: backend2.Spec{
				Style: id.FromString("litmus"),
			},
		},
		Compilers: compiler.MockSet(),
		Corpus:    corpus.Mock(),
	}
}
