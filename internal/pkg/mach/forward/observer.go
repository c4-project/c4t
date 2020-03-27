// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package forward

import (
	"encoding/json"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"
)

// Observer wraps a JSON encoder, lifting it to an Observer that sends JSON-encoded Forwards.
type Observer struct {
	*json.Encoder
}

// OnBuildStart sends a 'started' message through this Observer's encoder.
func (o *Observer) OnBuildStart(m builder.Manifest) {
	o.forwardHandlingError(Forward{BuildStart: &m})
}

// OnBuildRequest forwards r to this Observer's encoder.
func (o *Observer) OnBuildRequest(r builder.Request) {
	o.forwardHandlingError(Forward{BuildUpdate: &r})
}

// OnBuildRequest sends a 'finished' message through this Observer's encoder.
func (o *Observer) OnBuildFinish() {
	o.forwardHandlingError(Forward{BuildEnd: true})
}

// Error forwards err to this Observer's encoder.
func (o *Observer) Error(err error) {
	_ = o.forward(Forward{Error: err.Error()})
}

func (o *Observer) forwardHandlingError(f Forward) {
	if err := o.forward(f); err != nil {
		o.Error(err)
	}
}

func (o *Observer) forward(f Forward) error {
	return o.Encode(f)
}
