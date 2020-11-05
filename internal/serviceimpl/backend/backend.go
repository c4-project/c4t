// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package backend contains style-to-backend resolution.
package backend

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/stage/lifter"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
)

// Backend contains the various interfaces that a backend can implement.
type Backend interface {
	// Capabilities gets the capability set reported for this backend.
	Capabilities(b *backend.Spec) Capability

	// Some backends can lift test-cases into recipes (capability CanLift).
	lifter.SingleLifter

	// Backends that can be run standalone or produce executables (capability CanRunStandalone | CanProduceExecutables)
	// must give an observation parser for interpreting their stdout as observations.
	runner.ObsParser
}

// ErrNotSupported is the error that backends should return if we try to do something they don't support.
var ErrNotSupported = errors.New("service doesn't support action")
