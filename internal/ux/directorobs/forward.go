// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/director"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"
)

// Forward contains a director observation that has been forwarded from an instance to a 'main' observer,
// and disambiguating information.
type Forward struct {
	// Cycle is the cycle message if the forward kind is ForwardCycle;
	// if not, only its cycle is defined, and it determines the cycle from which this forward originates.
	Cycle director.CycleMessage

	// Kind is the kind of message that has been forwarded.
	Kind ForwardKind

	// Analysis is, when Kind is ForwardAnalysis, an analysis message.
	Analysis *analysis.Analysis

	// Compiler is, when Kind is ForwardCompilers, a forwarded compiler message.
	Compiler *compiler.Message

	// Save is, when Kind is ForwardSave, a forwarded save message.
	Save *saver.ArchiveMessage
}

type ForwardKind uint8

const (
	// ForwardCycle delimits a forwarding message where Cycle is populated.
	ForwardCycle ForwardKind = iota
	// ForwardAnalysis delimits a forwarding message where Analysis is populated (Cycle contains the cycle structure).
	ForwardAnalysis
	// ForwardCompiler delimits a forwarding message where Compiler is populated (Cycle contains the cycle structure).
	ForwardCompiler
	// ForwardSave delimits a forwarding message where Save is populated (Cycle contains the cycle structure).
	ForwardSave
)

// ForwardCycleMessage constructs a forwarding message for m.
func ForwardCycleMessage(m director.CycleMessage) Forward {
	return Forward{Cycle: m, Kind: ForwardCycle}
}

// ForwardAnalysisMessage constructs a forwarding message for m over cycle c.
func ForwardAnalysisMessage(c director.Cycle, m analysis.Analysis) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardAnalysis, Analysis: &m}
}

// ForwardCompilerMessage constructs a forwarding message for m over cycle c.
func ForwardCompilerMessage(c director.Cycle, m compiler.Message) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardCompiler, Compiler: &m}
}

// ForwardSaveMessage constructs a forwarding message for m over cycle c.
func ForwardSaveMessage(c director.Cycle, m saver.ArchiveMessage) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardSave, Save: &m}
}

// ForwardReceiver holds receive channels for Forward messages.
type ForwardReceiver observing.FanIn

// NewForwardReceiver constructs a new ForwardReceiver.
func NewForwardReceiver(f func(m Forward) error, cap int) *ForwardReceiver {
	return (*ForwardReceiver)(observing.NewFanIn(func(_ int, input interface{}) error {
		return f(input.(Forward))
	}, cap))
}

// Add adds a channel to the forward receiver.
func (r *ForwardReceiver) Add(c <-chan Forward) {
	(*observing.FanIn)(r).Add(c)
}

// Run runs the forward receiver.
func (r *ForwardReceiver) Run(ctx context.Context) error {
	return (*observing.FanIn)(r).Run(ctx)
}
