// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation

import (
	"context"
	"errors"
	"time"
)

var ErrNotActive = errors.New("automation is disabled in this config")

// Automator handles most of the legwork of automating mutant selection.
type Automator struct {
	// TickerF is a stubbable function used to create a ticker.
	TickerF func(time.Duration) (<-chan time.Time, Ticker)

	// config is the configuration that this automator accepts.
	config AutoConfig

	// mutantCh is the channel used to send mutants to the director.
	mutantCh chan Mutant

	// killCh is the channel used to receive kill signals from an observer.
	killCh chan Mutant

	// tickCh is the channel used to receive signals from the ticker.
	tickCh <-chan time.Time

	// ticker is the ticker used to handle time-slices.
	ticker Ticker

	// pool contains the current mutant set and handles which one is to be selected next.
	pool AutoPool
}

// NewAutomator constructs a new Automator given configuration cfg.
func NewAutomator(cfg AutoConfig) (*Automator, error) {
	if !cfg.IsActive() {
		return nil, ErrNotActive
	}

	a := &Automator{
		TickerF: StandardTicker,
		// killCh is lazily constructed by KillCh.
		config:   cfg,
		mutantCh: make(chan Mutant),
	}
	a.pool.Init(cfg.Mutants())
	return a, nil
}

// MutantCh gets a receive channel for taking mutants from this automator.
func (a *Automator) MutantCh() <-chan Mutant {
	return a.mutantCh
}

// KillCh gets a send channel for sending kill signals to this automator.
// If the automator isn't receiving kill signals, this will be nil.
// This channel must be closed.
func (a *Automator) KillCh() chan<- Mutant {
	if a.killCh == nil && a.config.ChangeKilled {
		a.killCh = make(chan Mutant)
	}
	return a.killCh
}

// Run runs this automator until ctx closes.
func (a *Automator) Run(ctx context.Context) {
	defer close(a.mutantCh)

	// Start off by sending the mutant at index 0.
	a.sendMutant(ctx)

	a.tickCh, a.ticker = a.TickerF(a.tickerDuration())
	defer a.ticker.Stop()

	a.loop(ctx)
}

func (a *Automator) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.drainKill()
			return
		case m := <-a.killCh:
			a.handleKill(ctx, m)
		case <-a.tickCh:
			a.handleTimeout(ctx)
		}
	}
}

func (a *Automator) handleKill(ctx context.Context, m Mutant) {
	a.pool.Kill(m)
	a.sendMutant(ctx)
	a.resetTicker()
}

func (a *Automator) handleTimeout(ctx context.Context) {
	a.pool.Advance()
	a.sendMutant(ctx)
}

// resetTicker resets the ticker after receiving a kill.
//
// This is important because, otherwise, the new mutant would only have the remainder of the old mutant's timeslice.
func (a *Automator) resetTicker() {
	a.ticker.Reset(a.tickerDuration())

	// A tick might have happened before we reset the ticker, in which case we eat it.
	// This isn't very robust.
	select {
	case <-a.tickCh:
	default:
	}
}

func (a *Automator) tickerDuration() time.Duration {
	return time.Duration(a.config.ChangeAfter)
}

// drainKill spins on any kill channel until it closes.
func (a *Automator) drainKill() {
	if a.killCh != nil {
		for range a.killCh {
		}
	}
}

func (a *Automator) sendMutant(ctx context.Context) {
	select {
	case <-ctx.Done():
	case a.mutantCh <- a.pool.Mutant():
	}
}

// Ticker is a mockable interface for time.Ticker.
type Ticker interface {
	Reset(d time.Duration)
	Stop()
}

//go:generate mockery --name Ticker

// StandardTicker gets a standard Go ticker if d is nonzero, and a no-op otherwise.
// In the latter case, the returned channel is nil.
func StandardTicker(d time.Duration) (<-chan time.Time, Ticker) {
	if d == 0 {
		return nil, nopTicker{}
	}

	t := time.NewTicker(d)
	return t.C, t
}

// nopTicker pretends to be a ticker, but never ticks.
type nopTicker struct{}

// Reset does nothing.
func (nopTicker) Reset(time.Duration) {}

// Stop does nothing.
func (nopTicker) Stop() {}

// AutoPool manages the next mutant to select in an automated mutation testing campaign.
//
// The policy AutoPool implements is:
// 1) Start by considering every mutant in turn.
// 2) Whenever a mutant is killed or its timeslot ends, advance to the next mutant.
// 3) If we are out of mutants, start again with the list of all mutants not killed by the steps above, and repeat.
// 4) If we kill every mutant, start again with every mutant.  (This behaviour may change eventually.)
type AutoPool struct {
	orig, curr, next []Mutant
	i                int
}

// Mutant gets the currently selected mutant.
func (a *AutoPool) Mutant() Mutant {
	return a.curr[a.i]
}

// Init initialises this pool with initial mutants muts.
func (a *AutoPool) Init(muts []Mutant) {
	a.orig = muts
	a.endCycle()
}

// Advance advances to the next mutant without killing it.
func (a *AutoPool) Advance() {
	a.next = append(a.next, a.Mutant())
	a.increment()
}

// Kill marks the mutant m as killed, if it is the current mutant.
func (a *AutoPool) Kill(m Mutant) {
	if a.Mutant().Index == m.Index {
		a.increment()
	}
}

func (a *AutoPool) increment() {
	a.i++
	if len(a.curr) <= a.i {
		a.endCycle()
	}
}

func (a *AutoPool) endCycle() {
	a.i = 0
	// Did we kill all of the mutants?  If so, we don't have any special support for this yet, so just reset entirely.
	if len(a.next) == 0 {
		a.next = a.orig
	}
	a.curr, a.next = a.next, make([]Mutant, 0, len(a.curr))
}
