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

	"github.com/mum4k/termdash/keyboard"

	"github.com/mum4k/termdash/cell"

	"github.com/MattWindsor91/act-tester/internal/machine"

	"github.com/mum4k/termdash/linestyle"

	"github.com/MattWindsor91/act-tester/internal/director/observer"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/mum4k/termdash/widgets/text"

	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"

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

	// obs contains the observer records used to populate machines.
	obs []*Observer
}

const (
	// idMachines is the container ID used to update the machine grid.
	idMachines = "machines"

	// MaxLogLines is the maximum number of lines that can be written to the log before it resets.
	MaxLogLines = 1000
)

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
func New() (*Dash, error) {
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

	logs := makeLogPane(d)

	c, err := container.New(d.term,
		container.SplitVertical(
			container.Left(logs),
			container.Right(container.ID("machines")),
			container.SplitPercent(25),
		),
	)
	if err != nil {
		return nil, err
	}

	d.container = c
	return &d, nil
}

func (d *Dash) OnMachines(m machine.Message) {
	switch m.Kind {
	case machine.MessageStart:
		d.setupMachineSplit(m.Index)
	case machine.MessageRecord:
		d.setupMachineID(m.Index, m.Machine.ID)
	}
}

func (d *Dash) logError(err error) {
	_ = d.log.Write(err.Error(), text.WriteCellOpts(cell.FgColor(cell.ColorRed)))
}

func makeLogPane(d Dash) container.Option {
	return container.SplitHorizontal(
		container.Top(
			container.SplitHorizontal(
				container.Top(
					container.Border(linestyle.Double), container.BorderTitle("Experiment Start"),
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

// machineContainerID calculates the container ID of the machine at location i.
// This is used to rename the container once we know its ID.
func machineContainerID(i int) string {
	return fmt.Sprintf("Machine%d", i)
}

func machineGridPercent(nmachines int) int {
	pc := 100 / nmachines
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
	return termdash.Run(ctx, d.term, d.container, termdash.KeyboardSubscriber(func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}))
}

// Close closes the dashboard's terminal.
func (d *Dash) Close() error {
	d.term.Close()
	return nil
}

// Instance locates the observer for the machine with ID mid.
func (d *Dash) Instance(mid id.ID) (observer.Instance, error) {
	o := d.machines[mid.String()]
	if o == nil {
		return nil, fmt.Errorf("instance not prepared for machine %s", mid.String())
	}
	return o, nil
}
