// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs

import (
	"context"
	"errors"

	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/stage/mach/observer"
	"github.com/c4-project/c4t/internal/stage/perturber"
	"github.com/c4-project/c4t/internal/stage/planner"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/observing"

	"github.com/c4-project/c4t/internal/director"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/stage/analyser/saver"
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

	// Instance is, when Kind is ForwardInstance, a forwarded instance message.
	Instance *director.InstanceMessage

	// Analysis is, when Kind is ForwardAnalysis, an analysis message.
	Analysis *analysis.Analysis

	// Compiler is, when Kind is ForwardCompilers, a forwarded compiler message.
	Compiler *compiler.Message

	// Save is, when Kind is ForwardSave, a forwarded save message.
	Save *saver.ArchiveMessage

	// Build is, when Kind is ForwardBuild, a forwarded build message.
	Build *builder.Message

	// Copy is, when Kind is ForwardCopy, a forwarded copy message.
	Copy *copier.Message
}

// ForwardKind is the enumeration of possible Forward messages.
type ForwardKind uint8

const (
	// ForwardCycle delimits a forwarding message where Cycle is populated.
	ForwardCycle ForwardKind = iota
	// ForwardInstance delimits a forwarding message where Instance is populated (Cycle contains the cycle structure).
	ForwardInstance
	// ForwardAnalysis delimits a forwarding message where Analysis is populated (Cycle contains the cycle structure).
	ForwardAnalysis
	// ForwardCompiler delimits a forwarding message where Compiler is populated (Cycle contains the cycle structure).
	ForwardCompiler
	// ForwardSave delimits a forwarding message where Save is populated (Cycle contains the cycle structure).
	ForwardSave
	// ForwardBuild delimits a forwarding message where Build is populated (Cycle contains the cycle structure).
	ForwardBuild
	// ForwardCopy delimits a forwarding message where Copy is populated (Cycle contains the cycle structure).
	ForwardCopy
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

// ForwardBuildMessage constructs a forwarding message for m over cycle c.
func ForwardBuildMessage(c director.Cycle, m builder.Message) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardBuild, Build: &m}
}

// ForwardCopyMessage constructs a forwarding message for m over cycle c.
func ForwardCopyMessage(c director.Cycle, m copier.Message) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardCopy, Copy: &m}
}

// ForwardInstanceMessage constructs a forwarding message for m over cycle c.
func ForwardInstanceMessage(c director.Cycle, m director.InstanceMessage) Forward {
	return Forward{Cycle: director.CycleMessage{Cycle: c}, Kind: ForwardInstance, Instance: &m}
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

	// OnCycleInstance handles a message for the instance of a particular cycle.
	OnCycleInstance(director.Cycle, director.InstanceMessage)

	// OnCycleAnalysis should handle an analysis for a particular cycle.
	OnCycleAnalysis(director.CycleAnalysis)

	// OnCycleBuild should handle a corpus build for a particular cycle.
	OnCycleBuild(director.Cycle, builder.Message)

	// OnCycleCompiler should handle a compiler message for a particular cycle.
	OnCycleCompiler(director.Cycle, compiler.Message)

	// OnCycleCopy should handle a copy message for a particular cycle.
	OnCycleCopy(director.Cycle, copier.Message)

	// OnCycleSave should handle an archive message for a particular cycle.
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
	case ForwardInstance:
		for _, h := range f.handlers {
			h.OnCycleInstance(fwd.Cycle.Cycle, *fwd.Instance)
		}
	case ForwardAnalysis:
		ca := director.CycleAnalysis{Cycle: fwd.Cycle.Cycle, Analysis: *fwd.Analysis}
		for _, h := range f.handlers {
			h.OnCycleAnalysis(ca)
		}
	case ForwardBuild:
		for _, h := range f.handlers {
			h.OnCycleBuild(fwd.Cycle.Cycle, *fwd.Build)
		}
	case ForwardCompiler:
		c := *fwd.Compiler
		for _, h := range f.handlers {
			h.OnCycleCompiler(fwd.Cycle.Cycle, c)
		}
	case ForwardCopy:
		for _, h := range f.handlers {
			h.OnCycleCopy(fwd.Cycle.Cycle, *fwd.Copy)
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

// OnAnalysis forwards a cycle analysis.
func (l *ForwardingInstanceObserver) OnAnalysis(c analysis.Analysis) {
	l.forward(ForwardAnalysisMessage(l.cycle, c))
}

// OnArchive forwards a cycle save message.
func (l *ForwardingInstanceObserver) OnArchive(s saver.ArchiveMessage) {
	l.forward(ForwardSaveMessage(l.cycle, s))
}

// OnPerturb does nothing, at the moment.
func (l *ForwardingInstanceObserver) OnPerturb(perturber.Message) {}

// OnPlan does nothing, at the moment.
func (l *ForwardingInstanceObserver) OnPlan(planner.Message) {}

// OnBuild forwards a cycle build message.
func (l *ForwardingInstanceObserver) OnBuild(m builder.Message) {
	l.forward(ForwardBuildMessage(l.cycle, m))
}

// OnCopy forwards a cycle copy message.
func (l *ForwardingInstanceObserver) OnCopy(m copier.Message) {
	l.forward(ForwardCopyMessage(l.cycle, m))
}

// OnMachineNodeAction does nothing.
func (l *ForwardingInstanceObserver) OnMachineNodeAction(observer.Message) {}

// OnInstance forwards an instance message, and closes the forwarding channel if the instance has closed.
func (l *ForwardingInstanceObserver) OnInstance(m director.InstanceMessage) {
	// TODO(@MattWindsor91): forward a message about this?
	l.forward(ForwardInstanceMessage(l.cycle, m))
	if m.Kind == director.KindInstanceClosed {
		close(l.fwd)
	}
}

func (l *ForwardingInstanceObserver) forward(f Forward) {
	select {
	case <-l.done:
	case l.fwd <- f:
	}
}
