// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"context"
	"errors"

	"github.com/MattWindsor91/c4t/internal/machine"
	"github.com/MattWindsor91/c4t/internal/model/id"

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
// alongside disambiguating information.
//
// This struct, and its sibling structs and interfaces,exist to solve the problem that many things in the director
// happen in instance threads, but then need to be observed by a single-threaded observer.
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

// ForwardKind is the enumeration of possible Forward messages.
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
//
// It is effectively a more type-safe form of observing.FanIn, and will usually be used inside a ForwardObserver.
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

// ForwardHandler is the interface of observers that can handle observations forwarded from an instance.
//
// These inject behaviour into a ForwardObserver.
type ForwardHandler interface {
	// ForwardHandler instances can observe cycles in the same way that the forwarding observer can.
	director.CycleObserver

	// These exist to let ForwardHandlers also implement bits of normal Observers, and can safely be zeroed.
	// TODO(@MattWindsor91): simplify this away.

	// ForwardHandler instances can observe machines.
	machine.Observer

	// ForwardHandler instances can observe preparations.
	director.PrepareObserver

	// OnCycleAnalysis should handle an analysis for a particular cycle.
	OnCycleAnalysis(director.CycleAnalysis)

	// OnCycleCompiler handles a compiler message for a particular cycle.
	OnCycleCompiler(director.Cycle, compiler.Message)

	// OnCycleSave handles an archive message for a particular cycle.
	OnCycleSave(director.Cycle, saver.ArchiveMessage)
}

// ForwardObserver is an observer that uses a ForwardReceiver and a ForwardHandler to handle observations.
type ForwardObserver struct {
	// done is the channel that will be closed when the observer finishes.
	done chan struct{}

	// handlers contain the observer logic.
	handlers []ForwardHandler
	// receiver contains the receiver, whose receiving loop should be spun up using Receiver.Run.
	receiver *ForwardReceiver
}

// OnMachines delegates to the forward handlers.
func (f *ForwardObserver) OnMachines(m machine.Message) {
	for _, h := range f.handlers {
		h.OnMachines(m)
	}
}

// OnCompilerConfig does nothing, for now.
func (f *ForwardObserver) OnCompilerConfig(compiler.Message) {
}

// OnBuild does nothing, for now.
func (f *ForwardObserver) OnBuild(builder.Message) {
}

// OnPlan does nothing, for now.
func (f *ForwardObserver) OnPlan(planner.Message) {
}

// OnPrepare forwards prepare messages to the handlers.
func (f *ForwardObserver) OnPrepare(p director.PrepareMessage) {
	for _, h := range f.handlers {
		h.OnPrepare(p)
	}
}

// ErrForwardHandlerNil occurs when NewForwardObserver is given a nil ForwardHandler.
var ErrForwardHandlerNil = errors.New("forward handler is nil")

// NewForwardObserver constructs a new ForwardObserver with the given handlers hs and message buffer capacity cap.
func NewForwardObserver(cap int, hs ...ForwardHandler) (*ForwardObserver, error) {
	for _, h := range hs {
		if h == nil {
			return nil, ErrForwardHandlerNil
		}
	}

	fo := &ForwardObserver{handlers: hs, done: make(chan struct{})}
	fo.receiver = NewForwardReceiver(fo.runStep, cap)
	return fo, nil
}

// Run runs the observer's forwarding loop,and closes any attached instances when it finishes.
func (f *ForwardObserver) Run(ctx context.Context) error {
	defer close(f.done)
	return f.receiver.Run(ctx)
}

func (f *ForwardObserver) runStep(fwd Forward) error {
	switch fwd.Kind {
	case ForwardAnalysis:
		ca := director.CycleAnalysis{Cycle: fwd.Cycle.Cycle, Analysis: *fwd.Analysis}
		for _, h := range f.handlers {
			h.OnCycleAnalysis(ca)
		}
	case ForwardCompiler:
		c := *fwd.Compiler
		for _, h := range f.handlers {
			h.OnCycleCompiler(fwd.Cycle.Cycle, c)
		}
	case ForwardCycle:
		for _, h := range f.handlers {
			h.OnCycle(fwd.Cycle)
		}
	case ForwardSave:
		s := *fwd.Save
		for _, h := range f.handlers {
			h.OnCycleSave(fwd.Cycle.Cycle, s)
		}
	}
	return nil
}

// Instance creates an instance observer that forwards to this logger.
func (f *ForwardObserver) Instance(id.ID) (director.InstanceObserver, error) {
	ch := make(chan Forward)
	f.receiver.Add(ch)
	return &ForwardingInstanceObserver{done: f.done, fwd: ch}, nil
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
