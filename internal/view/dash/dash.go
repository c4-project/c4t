// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package dash contains the act-tester console dashboard.
package dash

import (
	"context"
	"fmt"
	"time"

	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/mum4k/termdash/widgets/text"

	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/mum4k/termdash/container/grid"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
)

// Dash is a director observer that displays all of the current director machines in a terminal dashboard.
type Dash struct {
	container *container.Container
	term      terminalapi.Terminal
	startTime *text.Text
	log       *text.Text
	resultLog *ResultLog

	// nlines is the number of lines written, so far, to the log.
	nlines uint

	// machines maps from stringified machine IDs to their observers.
	// (There is a display order for the machines, but we don't track it ourselves.)
	machines map[string]*Observer
}

// MaxLogLines is the maximum number of lines that can be written to the log before it resets.
const MaxLogLines = 1000

// Write lets one write to the text console in the dash as if it were stderr.
func (d *Dash) Write(p []byte) (n int, err error) {
	// TODO(@MattWindsor91): mitigate against invalid input
	sp := string(p)

	d.nlines += countNewlines(sp)
	if MaxLogLines <= d.nlines {
		d.nlines -= MaxLogLines
		d.log.Reset()
		_ = d.log.Write("[log reset]")
	}

	return len(sp), d.log.Write(sp)
}

func countNewlines(sp string) uint {
	var n uint
	for _, c := range sp {
		if c == '\n' {
			n++
		}
	}
	return n
}

// New constructs a dashboard for the given machine IDs.
func New(mids []id.ID) (*Dash, error) {
	var (
		d   Dash
		err error
	)

	if d.term, err = tcell.New(); err != nil {
		return nil, err
	}

	if d.startTime, err = text.New(text.DisableScrolling()); err != nil {
		return nil, err
	}
	if err := d.startTime.Write(time.Now().Format(time.Stamp)); err != nil {
		return nil, err
	}

	if d.log, err = text.New(text.RollContent(), text.WrapAtWords()); err != nil {
		return nil, err
	}

	if d.resultLog, err = NewResultLog(); err != nil {
		return nil, err
	}

	var g []container.Option
	if d.machines, g, err = makeMachineGrid(mids, d.resultLog); err != nil {
		return nil, err
	}

	logs := makeLogPane(d)

	c, err := container.New(d.term,
		container.SplitVertical(
			container.Left(logs),
			container.Right(g...),
			container.SplitPercent(25),
		),
	)
	if err != nil {
		return nil, err
	}

	d.container = c
	return &d, nil
}

func makeLogPane(d Dash) container.Option {
	return container.SplitHorizontal(
		container.Top(
			container.SplitHorizontal(
				container.Top(
					container.Border(linestyle.Double), container.BorderTitle("Run Start"),
					container.PlaceWidget(d.startTime),
				),
				container.Bottom(
					container.Border(linestyle.Double), container.BorderTitle("System Log"), container.PlaceWidget(d.log),
				),
				container.SplitFixed(3),
			),
		),
		container.Bottom(
			container.Border(linestyle.Double), container.BorderTitle("Results Log"), container.PlaceWidget(d.resultLog.log),
		),
	)
}

func makeMachineGrid(mids []id.ID, rl *ResultLog) (map[string]*Observer, []container.Option, error) {
	gb := grid.New()

	obs := make(map[string]*Observer, len(mids))
	pc := machineGridPercent(mids)
	for _, mid := range mids {
		mstr := mid.String()
		var err error
		if obs[mstr], err = NewObserver(mid, rl); err != nil {
			return nil, nil, err
		}
		obs[mstr].AddToGrid(gb, mstr, pc)
	}

	g, err := gb.Build()
	return obs, g, err
}

func machineGridPercent(mids []id.ID) int {
	pc := 100 / len(mids)
	if pc == 100 {
		pc = 99
	}
	if pc == 0 {
		pc = 1
	}
	return pc
}

// Run runs a dashboard in a blocking manner.
func (d *Dash) Run(ctx context.Context, cancel func()) error {
	err := termdash.Run(ctx, d.term, d.container, termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyCtrlQ {
			cancel()
		}
	}))
	d.term.Close()
	return err
}

// Instance locates the observer for the machine with ID mid.
func (d *Dash) Instance(mid id.ID) (observer.Instance, error) {
	o := d.machines[mid.String()]
	if o == nil {
		return nil, fmt.Errorf("instance not prepared for machine %s", mid.String())
	}
	return o, nil
}
