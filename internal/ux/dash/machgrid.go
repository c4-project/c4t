// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
)

func (d *Dash) setupMachineSplit(nmachines int) {
	d.machines = make(map[string]*Observer, nmachines)
	if nmachines <= 0 {
		return
	}

	if err := d.updateMachineGrid(nmachines); err != nil {
		d.logError(err)
	}
}

func (d *Dash) setupMachineID(i int, mid id.ID) {
	midstr := mid.String()
	if err := d.container.Update(
		machineContainerID(i),
		container.BorderTitle(midstr),
	); err != nil {
		d.logError(err)
	}
	d.machines[midstr] = d.obs[i]
}

// updateMachineGrid populates the dash's observer list and machine grid with nmachines machine observers.
func (d *Dash) updateMachineGrid(nmachines int) error {
	var (
		g   []container.Option
		err error
	)
	if d.obs, g, err = makeMachineGrid(nmachines, d.resultLog); err != nil {
		return err
	}
	return d.container.Update(idMachines, g...)
}

func makeMachineGrid(nmachines int, rl *ResultLog) ([]*Observer, []container.Option, error) {
	gb := grid.New()

	obs := make([]*Observer, nmachines)
	pc := machineGridPercent(nmachines)
	for i := range obs {
		var err error
		if obs[i], err = NewObserver(rl); err != nil {
			return nil, nil, err
		}
		obs[i].AddToGrid(gb, machineContainerID(i), pc)
	}

	g, err := gb.Build()
	return obs, g, err
}
