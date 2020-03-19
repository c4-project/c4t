// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package dash contains the act-tester console dashboard.
package dash

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/mum4k/termdash/widgets/text"

	"github.com/mum4k/termdash/linestyle"

	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/mum4k/termdash/widgets/gauge"

	"github.com/mum4k/termdash/container/grid"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
)

// Dash is a director observer that displays all of the current director machines in a terminal dashboard.
type Dash struct {
	container *container.Container
	term      terminalapi.Terminal
	log       *text.Text

	// machines maps from stringified machine IDs to their observers.
	// (There is a display order for the machines, but we don't track it ourselves.)
	machines map[string]*Observer
}

// Write lets one write to the text console in the dash as if it were stderr.
func (d *Dash) Write(p []byte) (n int, err error) {
	// TODO(@MattWindsor91): mitigate against invalid input
	sp := string(p)
	return len(sp), d.log.Write(sp)
}

// New constructs a dashboard for the given machine IDs.
func New(mids []id.ID) (*Dash, error) {
	t, err := termbox.New()
	if err != nil {
		return nil, err
	}

	x, err := text.New()
	if err != nil {
		return nil, err
	}

	obs, g, err := makeMachineGrid(mids)
	if err != nil {
		return nil, err
	}
	g = append(g, container.Border(linestyle.Light))
	c, err := container.New(t,
		container.SplitVertical(
			container.Left(container.Border(linestyle.Double), container.BorderTitle("Log"), container.PlaceWidget(x)),
			container.Right(g...),
		),
	)
	if err != nil {
		return nil, err
	}

	d := Dash{
		container: c,
		term:      t,
		log:       x,
		machines:  obs,
	}
	return &d, nil
}

func makeMachineGrid(mids []id.ID) (map[string]*Observer, []container.Option, error) {
	gb := grid.New()

	obs := make(map[string]*Observer, len(mids))
	pc := machineGridPercent(mids)
	for _, mid := range mids {
		mstr := mid.String()
		var err error
		if obs[mstr], err = makeMachine(); err != nil {
			return nil, nil, err
		}
		addMachineToGrid(gb, mstr, obs[mstr], pc)
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

func makeMachine() (*Observer, error) {
	var d Observer

	var err error

	if d.g, err = gauge.New(); err != nil {
		return nil, err
	}

	if d.last, err = text.New(text.RollContent()); err != nil {
		return nil, err
	}

	return &d, nil
}

func addMachineToGrid(gb *grid.Builder, midstr string, o *Observer, pc int) {
	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{container.Border(linestyle.Double), container.BorderTitle(midstr)},
		grid.RowHeightFixed(1, grid.Widget(o.g)),
		grid.RowHeightFixed(1, grid.Widget(o.last)),
	))
}

// Run runs a dashboard in a blocking manner.
func (d *Dash) Run(ctx context.Context, cancel func()) error {
	err := termdash.Run(ctx, d.term, d.container, termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}))
	d.term.Close()
	return err
}

// Machine locates the observer for the machine with ID mid.
func (d *Dash) Machine(mid id.ID) builder.Observer {
	o := d.machines[mid.String()]
	if o == nil {
		return builder.SilentObserver{}
	}
	return o
}
