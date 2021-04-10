// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package observing

import (
	"context"
	"reflect"
)

const (
	// doneCh is the index of the 'done' channel in a fan-in.
	doneCh = 0
	// nReservedCh is number of channels reserved at the top of a fan-in.
	nReservedCh = doneCh + 1
)

// FanIn is a low-level, reflection-based device for forwarding observations from multiple concurrent sources to a
// single handler.
type FanIn struct {
	f      func(i int, input interface{}) error
	err    error
	chans  []reflect.SelectCase
	nchans int
}

// NewFanIn creates a fan-in with the given handling function f and capacity hint cap.
func NewFanIn(f func(i int, input interface{}) error, cap int) *FanIn {
	return &FanIn{
		f:      f,
		chans:  make([]reflect.SelectCase, nReservedCh, cap+nReservedCh),
		nchans: 0,
		err:    nil,
	}
}

// Add adds a channel to a fan-in.
func (f *FanIn) Add(ch interface{}) {
	f.chans = append(f.chans, recvCase(ch))
	f.nchans++
}

// Run runs the fan-in on context ctx.
// It is not re-entrant.
func (f *FanIn) Run(ctx context.Context) error {
	f.chans[doneCh] = recvCase(ctx.Done())

	for {
		// Note that this lets us terminate *before* ctx is cancelled, if every non-ctx channel closes.
		if f.nchans == 0 {
			return f.err
		}

		chosen, recv, recvOK := reflect.Select(f.chans)
		switch {
		case chosen == doneCh:
			f.ctxDone(ctx)
		case !recvOK:
			f.remove(chosen)
		case f.err == nil:
			f.err = f.f(chosen-nReservedCh, recv.Interface())
		}
	}
}

func (f *FanIn) ctxDone(ctx context.Context) {
	f.err = ctx.Err()
	// We can't use remove here, as it'll decrement nchans.
	f.zeroCh(doneCh)
}

func (f *FanIn) remove(i int) {
	f.zeroCh(i)
	f.nchans--
}

func (f *FanIn) zeroCh(i int) {
	f.chans[i].Chan = reflect.Value{}
}

func recvCase(ch interface{}) reflect.SelectCase {
	return reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
}
