// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package forward

import (
	"encoding/json"
	"io"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/observer"

	"github.com/MattWindsor91/act-tester/internal/model/corpus/builder"
)

// Observer wraps a JSON encoder, lifting it to an Observer that sends JSON-encoded Forwards.
type Observer struct {
	*json.Encoder
}

// NewObserver creates a forwarding observer over w.
func NewObserver(w io.Writer) *Observer {
	return &Observer{json.NewEncoder(w)}
}

// OnBuild sends a build message through this Observer's encoder.
func (o *Observer) OnBuild(m builder.Message) {
	o.forwardHandlingError(Forward{Build: &m})
}

// OnMachineNodeAction sends an action message through this Observer's encoder.
func (o *Observer) OnMachineNodeAction(m observer.Message) {
	o.forwardHandlingError(Forward{Action: &m})
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
