// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the top-level ACT test director, which manages a full testing campaign.
package director

import (
	"context"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// config is the configuration for the director.
	config *Config

	// l is the logger for the director.
	l *log.Logger
}

// New creates a new Director given a global act-tester config.
// It fails if the config is missing or ill-formed.
func New(c *Config) (*Director, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	return &Director{config: c, l: iohelp.EnsureLog(c.Logger)}, nil
}

func checkConfig(c *Config) error {
	if c == nil {
		return ErrConfigNil
	}
	if c.Paths == nil {
		return iohelp.ErrPathsetNil
	}
	if c.Machines == nil || len(c.Machines) == 0 {
		return ErrNoMachines
	}
	return nil
}

// Direct runs the director d.
func (d *Director) Direct(ctx context.Context) error {
	d.l.Print("making directories")
	if err := d.config.Paths.Prepare(); err != nil {
		return err
	}

	eg, ectx := errgroup.WithContext(ctx)
	for midstr, c := range d.config.Machines {
		m, err := d.makeMachine(midstr, c)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			return m.Run(ectx)
		})
	}

	return eg.Wait()
}

func (d *Director) makeMachine(midstr string, c config.Machine) (*Machine, error) {
	l := log.New(d.l.Writer(), logPrefix(midstr), 0)
	mid, err := model.TryIDFromString(midstr)
	if err != nil {
		return nil, err
	}
	m := Machine{
		Config: c,
		ID:     mid,
		Paths:  d.config.Paths.MachineScratch(mid),
		Logger: l,
	}
	return &m, nil
}

func logPrefix(midstr string) string {
	return midstr + ": "
}
