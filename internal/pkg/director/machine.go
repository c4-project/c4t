// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Machine contains the state necessary to run a single machine loop of a director.
type Machine struct {
	// ID is the ID for this machine.
	ID model.ID

	// Logger points to a logger for this machine's loop.
	Logger *log.Logger

	// Config contains the machine config for this machine.
	Config config.Machine

	// Paths contains the scratch pathset for this machine.
	Paths *MachinePathset
}

// Run runs this machine's testing loop.
func (m *Machine) Run(_ context.Context) error {
	m.Logger = iohelp.EnsureLog(m.Logger)
	if m.Paths == nil {
		return fmt.Errorf("%w: paths for machine %s", iohelp.ErrPathsetNil, m.ID.String())
	}

	if m.Config.SSH != nil {
		return errors.New("TODO: SSH support not yet available")
	}

	m.Logger.Print("preparing scratch directories")
	if err := m.Paths.Prepare(); err != nil {
		return err
	}

	m.Logger.Print("starting loop")

	// TODO(@MattWindsor91)
	return nil
}
