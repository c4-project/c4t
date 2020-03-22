// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package director contains the top-level ACT test director, which manages a full testing campaign.
package director

import (
	"context"
	"fmt"
	"log"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"

	"golang.org/x/sync/errgroup"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

// Director contains the main state and configuration for the test director.
type Director struct {
	// config is the configuration for the director.
	config *Config

	// files is the input file set.
	files []string

	// l is the logger for the director.
	l *log.Logger
}

// New creates a new Director given a global act-tester config and the input file set files.
// It fails if the config is missing or ill-formed.
func New(c *Config, files []string) (*Director, error) {
	if len(files) == 0 {
		return nil, liftInitError(corpus.ErrNone)
	}

	if err := c.Check(); err != nil {
		return nil, liftInitError(err)
	}

	return &Director{config: c, files: files, l: iohelp.EnsureLog(c.Logger)}, nil
}

// liftInitError lifts err to mention that it occurred during initialisation of a director.
func liftInitError(err error) error {
	return fmt.Errorf("while initialising director: %w", err)
}

// Direct runs the director d.
func (d *Director) Direct(ctx context.Context) error {
	d.l.Println("making directories")
	if err := d.config.Paths.Prepare(); err != nil {
		return err
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ectx := errgroup.WithContext(cctx)
	for midstr, c := range d.config.Machines {
		m, err := d.makeMachine(midstr, c)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			return m.Run(ectx)
		})
	}

	eg.Go(func() error {
		return d.config.Observer.Run(ectx, cancel)
	})

	return eg.Wait()
}

func (d *Director) makeMachine(midstr string, c config.Machine) (*Instance, error) {
	l := log.New(d.l.Writer(), logPrefix(midstr), 0)
	mid, err := id.TryFromString(midstr)
	if err != nil {
		return nil, err
	}
	obs, err := d.config.Observer.Instance(mid)
	if err != nil {
		return nil, err
	}
	sps := d.config.Paths.MachineScratch(mid)
	vps := d.config.Paths.MachineSaved(mid)
	m := Instance{
		MachConfig:   c,
		SSHConfig:    d.config.SSH,
		Env:          &d.config.Env,
		ID:           mid,
		InFiles:      d.files,
		Observer:     obs,
		ScratchPaths: sps,
		SavedPaths:   vps,
		Logger:       l,
		Quantities:   d.config.Quantities,
	}
	return &m, nil
}

func logPrefix(midstr string) string {
	return midstr + ": "
}
