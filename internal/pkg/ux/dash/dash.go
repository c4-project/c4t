// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package dash contains the act-tester console dashboard.
package dash

import (
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/observer"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/mum4k/termdash/widgets/text"

	"github.com/mum4k/termdash/linestyle"

	"github.com/mum4k/termdash/terminal/terminalapi"

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

	x, err := text.New(text.RollContent())
	if err != nil {
		return nil, err
	}

	obs, g, err := makeMachineGrid(mids)
	if err != nil {
		return nil, err
	}
	c, err := container.New(t,
		container.SplitVertical(
			container.Left(container.Border(linestyle.Double), container.BorderTitle("Log"), container.PlaceWidget(x)),
			container.Right(g...),
			container.SplitPercent(30),
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
		if obs[mstr], err = NewObserver(); err != nil {
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
		if k.Key == 'q' || k.Key == 'Q' {
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
