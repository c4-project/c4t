// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package directorobs contains a pre-packaged observer set for the test director.
package directorobs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MattWindsor91/c4t/internal/director"

	"github.com/MattWindsor91/c4t/internal/config"
	"github.com/MattWindsor91/c4t/internal/helper/errhelp"
	"github.com/MattWindsor91/c4t/internal/ux/dash"
	"golang.org/x/sync/errgroup"
)

// Obs is the standard top-level director observer.
type Obs struct {
	// dash is the dashboard, if enabled (null otherwise).
	dash *dash.Dash
	// resultLog is a forward handler that logs results and other significant cycle phenomena to a text file.
	resultLog *Logger
	// statPersister is a forward handler that persists statistics in a JSON file.
	statPersister *StatPersister
	// fwd contains the forwarding observer that hosts forward handlers.
	fwd *ForwardObserver

	// TODO(@MattWindsor91): in the medium term, I expect director observers that aren't forwarding will disappear.
	// This will probably happen if/when the dashboard gets decoupled into a separate networked tool.
}

// NewObs creates a director observer using the global configuration cfg.
// If useDash is true, it will create a dashboard; otherwise, it will bypass this.
func NewObs(cfg *config.Config, useDash bool) (*Obs, error) {
	obs := new(Obs)
	if err := obs.setup(cfg, useDash); err != nil {
		_ = obs.Close()
		return nil, err
	}
	return obs, nil
}

func (o *Obs) setup(cfg *config.Config, useDash bool) error {
	var err error

	if o.resultLog, err = loggerFromConfig(cfg); err != nil {
		return fmt.Errorf("while creating logger: %w", err)
	}
	if o.statPersister, err = statPersisterFromConfig(cfg); err != nil {
		return fmt.Errorf("while creating stat persister: %w", err)
	}
	if err = o.setupForwarder(); err != nil {
		return err
	}
	if !useDash {
		return nil
	}
	return o.setupDash()
}

func (o *Obs) setupDash() error {
	var err error
	o.dash, err = dash.New()
	return err
}

func (o *Obs) setupForwarder() error {
	fhs := make([]ForwardHandler, 0, 2)
	if o.resultLog != nil {
		fhs = append(fhs, o.resultLog)
	}
	if o.statPersister != nil {
		fhs = append(fhs, o.statPersister)
	}
	// TODO(@MattWindsor91): wire cap up to number of instances
	var err error
	o.fwd, err = NewForwardObserver(10, fhs...)
	return err
}

// loggerFromConfig constructs a logger according to the configuration in cfg.
func loggerFromConfig(cfg *config.Config) (*Logger, error) {
	logw, err := createResultLogFile(cfg)
	if err != nil {
		return nil, err
	}
	return NewLogger(logw, log.LstdFlags)
}

func createResultLogFile(c *config.Config) (*os.File, error) {
	logpath, err := c.Paths.OutPath("results.log")
	if err != nil {
		return nil, fmt.Errorf("expanding result log file path: %w", err)
	}
	logw, err := os.Create(logpath)
	if err != nil {
		return nil, fmt.Errorf("opening result log file: %w", err)
	}
	return logw, nil
}

func statPersisterFromConfig(c *config.Config) (*StatPersister, error) {
	path, err := c.Paths.OutPath("stats.json")
	if err != nil {
		return nil, fmt.Errorf("expanding stat persister file path: %w", err)
	}
	f, err := OpenStatFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening stat persister file %q: %w", path, err)
	}
	return NewStatPersister(f)
}

func (o *Obs) Observers() []director.Observer {
	if o.dash != nil {
		return []director.Observer{o.dash, o.fwd}
	}
	return []director.Observer{o.fwd}
}

func (o *Obs) Run(ctx context.Context, cancel context.CancelFunc) error {
	eg, ectx := errgroup.WithContext(ctx)
	if o.dash != nil {
		eg.Go(func() error {
			return o.dash.Run(ectx, cancel)
		})
	}
	eg.Go(func() error {
		return o.fwd.Run(ectx)
	})
	return eg.Wait()
}

func (o *Obs) Close() error {
	var derr, rerr, serr error
	if o.dash != nil {
		derr = o.dash.Close()
	}
	if o.resultLog != nil {
		rerr = o.resultLog.Close()
	}
	if o.statPersister != nil {
		serr = o.statPersister.Close()
	}
	return errhelp.FirstError(derr, rerr, serr)
}
