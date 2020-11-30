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
	dash      *dash.Dash
	resultLog *Logger
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

	if o.resultLog, err = LoggerFromConfig(cfg); err != nil {
		return err
	}
	if !useDash {
		return nil
	}
	o.dash, err = dash.New()
	return err
}

// LoggerFromConfig constructs a logger according to the configuration in cfg.
func LoggerFromConfig(cfg *config.Config) (*Logger, error) {
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

func (o *Obs) Observers() []director.Observer {
	if o.dash == nil {
		return []director.Observer{o.resultLog}
	}
	return []director.Observer{o.dash, o.resultLog}
}

func (o *Obs) Run(ctx context.Context, cancel context.CancelFunc) error {
	eg, ectx := errgroup.WithContext(ctx)
	if o.dash != nil {
		eg.Go(func() error {
			return o.dash.Run(ectx, cancel)
		})
	}
	eg.Go(func() error {
		return o.resultLog.Run(ectx)
	})
	return eg.Wait()
}

func (o *Obs) Close() error {
	var derr, rerr error
	if o.dash != nil {
		derr = o.dash.Close()
	}
	if o.resultLog != nil {
		rerr = o.resultLog.Close()
	}
	return errhelp.FirstError(derr, rerr)
}
