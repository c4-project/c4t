// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"context"

	"github.com/MattWindsor91/c4t/internal/copier"
	"github.com/MattWindsor91/c4t/internal/stage/mach/observer"
	"github.com/MattWindsor91/c4t/internal/stage/perturber"
	"github.com/MattWindsor91/c4t/internal/stage/planner"
	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

	"github.com/MattWindsor91/c4t/internal/observing"

	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	"github.com/MattWindsor91/c4t/internal/plan/analysis"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"
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

// ForwardingInstanceObserver is an instance observer that just forwards every observation to a director observer.
type ForwardingInstanceObserver struct {
	// done is a channel closed when the instance can no longer log.
	done <-chan struct{}
	// fwd is a channel to which the instance forwards messages to the main logger.
	fwd chan<- Forward
	// cycle contains information about the current iteration.
	cycle director.Cycle
}

func (l *ForwardingInstanceObserver) OnCompilerConfig(m compiler.Message) {
	l.forward(ForwardCompilerMessage(l.cycle, m))
}

// OnIteration notes that the instance's iteration has changed.
func (l *ForwardingInstanceObserver) OnCycle(r director.CycleMessage) {
	if r.Kind == director.CycleStart {
		l.cycle = r.Cycle
	}
	l.forward(ForwardCycleMessage(r))
}

// OnCollation logs a collation to this logger.
func (l *ForwardingInstanceObserver) OnAnalysis(c analysis.Analysis) {
	l.forward(ForwardAnalysisMessage(l.cycle, c))
}

func (l *ForwardingInstanceObserver) OnArchive(s saver.ArchiveMessage) {
	l.forward(ForwardSaveMessage(l.cycle, s))
}

// OnPerturb does nothing, at the moment.
func (l *ForwardingInstanceObserver) OnPerturb(perturber.Message) {}

// OnPlan does nothing, at the moment.
func (l *ForwardingInstanceObserver) OnPlan(planner.Message) {}

// OnBuild does nothing.
func (l *ForwardingInstanceObserver) OnBuild(builder.Message) {}

// OnCopy does nothing.
func (l *ForwardingInstanceObserver) OnCopy(copier.Message) {}

// OnMachineNodeMessage does nothing.
func (l *ForwardingInstanceObserver) OnMachineNodeAction(observer.Message) {}

// OnInstanceClose doesn't (yet) log the instance closure, but closes the forwarding channel.
func (l *ForwardingInstanceObserver) OnInstanceClose() {
	// TODO(@MattWindsor91): forward a message about this?
	close(l.fwd)
}

func (l *ForwardingInstanceObserver) forward(f Forward) {
	select {
	case <-l.done:
	case l.fwd <- f:
	}
}
