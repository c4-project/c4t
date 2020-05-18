// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Mock makes a valid, but mocked-up, plan.
func Mock() *Plan {
	// TODO(@MattWindsor91): add things to this plan as time goes on.
	return &Plan{
		Metadata: Header{
			Creation: time.Date(2011, time.November, 11, 11, 11, 11, 0, time.FixedZone("PST", -8*60*60)),
			Seed:     8675309,
			Version:  CurrentVer,
		},
		Machine: NamedMachine{
			ID: id.FromString("localhost"),
		},
		Backend: &service.Backend{
			Style: id.FromString("litmus"),
		},
		Compilers: nil,
		Corpus:    corpus.Mock(),
	}
}
