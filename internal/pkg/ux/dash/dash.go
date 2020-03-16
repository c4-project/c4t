// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package dash contains the act-tester console dashboard.
package dash

import (
	"context"
	"reflect"
	"sort"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/mum4k/termdash/widgets/text"

	"github.com/mum4k/termdash/linestyle"

	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/mum4k/termdash/widgets/gauge"

	"github.com/mum4k/termdash/container/grid"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
)

// Dash is a director observer that displays all of the current director machines in a terminal dashboard.
type Dash struct {
	container *container.Container
	term      terminalapi.Terminal
	log       *text.Text
	machines  []Observer
}

// Write lets one write to the text console in the dash as if it were stderr.
func (d *Dash) Write(p []byte) (n int, err error) {
	// TODO(@MattWindsor91): mitigate against invalid input
	sp := string(p)
	return len(sp), d.log.Write(sp)
}

// New constructs a dashboard for the given machine IDs.
func New(mids []model.ID) (*Dash, error) {
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

func makeMachineGrid(mids []model.ID) ([]Observer, []container.Option, error) {
	gb := grid.New()

	obs := make([]Observer, len(mids))
	pc := 100 / len(mids)
	if pc == 100 {
		pc = 99
	}
	if pc == 0 {
		pc = 1
	}
	for i, mid := range mids {
		if err := addMachine(&obs[i], mid, gb, pc); err != nil {
			return nil, nil, err
		}
	}

	g, err := gb.Build()
	return obs, g, err
}

func addMachine(d *Observer, mid model.ID, gb *grid.Builder, pc int) error {
	d.mid = mid

	var err error

	if d.g, err = gauge.New(); err != nil {
		return err
	}

	if d.last, err = text.New(text.RollContent()); err != nil {
		return err
	}

	gb.Add(grid.RowHeightPercWithOpts(pc,
		[]container.Option{container.Border(linestyle.Double), container.BorderTitle(mid.String())},
		grid.RowHeightFixed(1, grid.Widget(d.g)),
		grid.RowHeightFixed(1, grid.Widget(d.last)),
	))

	return err
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
func (d *Dash) Machine(mid model.ID) builder.Observer {
	n := len(d.machines)
	i := sort.Search(n, func(i int) bool {
		return reflect.DeepEqual(mid, d.machines[i].mid)
	})
	if n <= i {
		return builder.SilentObserver{}
	}
	return &d.machines[i]
}
