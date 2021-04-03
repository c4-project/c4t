// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/c4-project/c4t/internal/id"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
)

// TODO(@MattWindsor91): can this be moved to Instance?
func (d *Dash) assignMachineID(i int, mid id.ID) {
	midstr := mid.String()
	if err := d.container.Update(
		machineContainerID(i),
		container.BorderTitle(midstr),
	); err != nil {
		d.logError(err)
	}
}

// updateMachineGrid updates the dash's instance grid to accommodate ninst instances.
func (d *Dash) updateMachineGrid() error {
	var (
		g   []container.Option
		err error
	)
	if g, err = d.makeMachineGrid(); err != nil {
		return err
	}
	return d.container.Update(idInstances, g...)
}

func (d *Dash) makeMachineGrid() ([]container.Option, error) {
	gb := grid.New()

	pc := machineGridPercent(len(d.instances))
	for _, in := range d.instances {
		in.AddToGrid(gb, pc)
	}

	return gb.Build()
}
