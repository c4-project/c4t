// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package dash contains the c4t console dashboard.
package dash

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MattWindsor91/c4t/internal/director"
	"github.com/MattWindsor91/c4t/internal/stage/analyser/saver"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"
	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

	"github.com/MattWindsor91/c4t/internal/stage/planner"

	"github.com/mum4k/termdash/keyboard"

	"github.com/mum4k/termdash/cell"

	"github.com/MattWindsor91/c4t/internal/machine"

	"github.com/mum4k/termdash/linestyle"

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

	// instances contains all instances currently allocated for the dashboard.
	instances []*Instance
}

// OnCycle forwards the cycle message m to the relevant instance.
func (d *Dash) OnCycle(m director.CycleMessage) {
	// TODO(@MattWindsor91): the instance itself should handle this
	if m.Kind == director.CycleStart {
		d.assignMachineID(m.Cycle.Instance, m.Cycle.MachineID)
	}
	d.onInstance(m.Cycle.Instance, func(i *Instance) { i.OnCycle(m) })
}

// OnCycleAnalysis forwards the analysis m to the relevant instance.
func (d *Dash) OnCycleAnalysis(m director.CycleAnalysis) {
	d.onInstance(m.Cycle.Instance, func(i *Instance) { i.OnAnalysis(m.Analysis) })
}

// OnCycleCompiler forwards the compiler message m to the instance mentioned in c.
func (d *Dash) OnCycleCompiler(c director.Cycle, m compiler.Message) {
	d.onInstance(c.Instance, func(i *Instance) { i.OnCompilerConfig(m) })
}

// OnCycleSave forwards the archive message m to the instance mentioned in c.
func (d *Dash) OnCycleSave(c director.Cycle, m saver.ArchiveMessage) {
	d.onInstance(c.Instance, func(i *Instance) { i.OnArchive(m) })
}

// ErrNoSuchInstance occurs when a message arrives from an instance that the dashboard hasn't allocated room for.
var ErrNoSuchInstance = errors.New("received message for instance that doesn't exist")

func (d *Dash) onInstance(i int, f func(*Instance)) {
	if i < 0 || len(d.instances) <= i {
		d.logError(fmt.Errorf("%w: %d", ErrNoSuchInstance, i))
		return
	}
	f(d.instances[i])
}

// ensureInstances ensures there are at least n instances, and recalculates the machine grid if not.
func (d *Dash) ensureInstances(n int) error {
	var err error

	ninst := len(d.instances)
	excess := n - ninst
	if excess <= 0 {
		return nil
	}
	newinsts := make([]*Instance, excess)
	for j := range newinsts {
		if newinsts[j], err = NewInstance(machineContainerID(j+ninst), d); err != nil {
			return err
		}
	}

	d.instances = append(d.instances, newinsts...)
	return d.updateMachineGrid()
}

const (
	// idInstances is the container ID used to update the instance grid.
	idInstances = "machines"

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
			container.Right(container.ID(idInstances)),
			container.SplitPercent(25),
		),
	)
	if err != nil {
		return nil, err
	}

	d.container = c
	return &d, nil
}

func (d *Dash) OnMachines(machine.Message) {
	// do nothing, for now
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

// OnPrepare uses the instance calculation to prepare a machine grid.
func (d *Dash) OnPrepare(m director.PrepareMessage) {
	// TODO(@MattWindsor91): broadcast the quantities somewhere
	if m.Kind != director.PrepareInstances {
		return
	}
	if err := d.log.Write("Instances: " + strconv.Itoa(m.NumInstances)); err != nil {
		d.logError(err)
	}
	if err := d.ensureInstances(m.NumInstances); err != nil {
		d.logError(err)
	}
}

// OnCompilerConfig (currently) does nothing.
func (d *Dash) OnCompilerConfig(compiler.Message) {
	// NB: this version of OnCompilerConfig gets triggered by planner messages;
	// the one in Instance is the one thar gets triggered by perturber messages.
}

// OnBuild (currently) does nothing.
func (d *Dash) OnBuild(builder.Message) {
	// NB: this version of OnBuild gets triggered by planner messages;
	// the one in Instance is the one thar gets triggered by perturber messages.
}

// OnPlan (currently) does nothing.
func (d *Dash) OnPlan(planner.Message) {
	// TODO(@MattWindsor91): do something here
}
