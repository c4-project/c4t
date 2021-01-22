// Copyright (c) 2021 Matt Windsor and contributors
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
	killCh chan struct{}

	// tickCh is the channel used to receive signals from the ticker.
	tickCh <-chan time.Time

	// ticker is the ticker used to handle time-slices.
	ticker Ticker

	// i is the current index in mutants.
	i int
	// mutants is the mutant list.
	mutants []Mutant
}

// NewAutomator constructs a new Automator given configuration cfg.
func NewAutomator(cfg AutoConfig) (*Automator, error) {
	if !cfg.IsActive() {
		return nil, ErrNotActive
	}

	return &Automator{
		TickerF: StandardTicker,
		// killCh is lazily constructed by KillCh.
		config:   cfg,
		mutantCh: make(chan Mutant),
		i:        0,
		mutants:  cfg.Mutants(),
	}, nil
}

// MutantCh gets a receive channel for taking mutants from this automator.
func (a *Automator) MutantCh() <-chan Mutant {
	return a.mutantCh
}

// KillCh gets a send channel for sending kill signals to this automator.
// If the automator isn't receiving kill signals, this will be nil.
// This channel must be closed.
func (a *Automator) KillCh() chan<- struct{} {
	if a.killCh == nil && a.config.ChangeKilled {
		a.killCh = make(chan struct{})
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
		case <-a.killCh:
			a.sendMutant(ctx)
			a.resetTicker()
		case <-a.tickCh:
			a.sendMutant(ctx)
		}
	}
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
	if a.killCh == nil {
		return
	}
	for range a.killCh {
	}
}

func (a *Automator) sendMutant(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case a.mutantCh <- a.mutants[a.i]:
		a.advanceIndex()
	}
}

// advanceIndex increments the mutant index, wrapping around at the end of the mutant list.
func (a *Automator) advanceIndex() {
	// NewAutomator checks that 0 < len(a.mutants).
	a.i = (a.i + 1) % len(a.mutants)
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
