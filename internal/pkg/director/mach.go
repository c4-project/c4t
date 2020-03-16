// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

	"github.com/MattWindsor91/act-tester/internal/pkg/mach/forward"

	"golang.org/x/sync/errgroup"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// LocalMach runs the 'mach' binary locally.
type LocalMach struct {
	// Dir is the directory in which we are running the machine-runner.
	Dir string

	// Observer is the observer to which we are sending updates from the machine-runner.
	Observer builder.Observer
}

// Run runs act-tester-mach locally.
func (m *LocalMach) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	eg, ectx := errgroup.WithContext(ctx)

	// TODO(@MattWindsor91): remote invocation
	// TODO(@MattWindsor91): observation
	cmd := exec.CommandContext(ectx, "act-tester-mach", "-J", "-d", m.Dir)
	stdin, stdout, stderr, err := openPipes(cmd)
	if err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		_ = stderr.Close()
		_ = stdout.Close()
		_ = stdin.Close()
		return nil, fmt.Errorf("while starting local runner: %w", err)
	}

	var p2 plan.Plan
	eg.Go(func() error {
		return sendPlan(p, stdin)
	})
	eg.Go(func() error {
		if _, err := toml.DecodeReader(stdout, &p2); err != nil {
			return fmt.Errorf("while decoding the output plan: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		r := forward.Replayer{
			Decoder: json.NewDecoder(stderr),
			Obs:     m.Observer,
		}
		return r.Run(ectx)
	})

	// Waiting _should_ close the pipes.
	err = eg.Wait()
	werr := cmd.Wait()

	if err != nil {
		return nil, err
	}
	return &p2, werr
}

// openPipes tries to open stdin, stdout, and stderr pipes for c.
func openPipes(c *exec.Cmd) (stdin io.WriteCloser, stdout, stderr io.ReadCloser, err error) {
	if stdin, err = c.StdinPipe(); err != nil {
		return nil, nil, nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	if stdout, err = c.StdoutPipe(); err != nil {
		_ = stdin.Close()
		return nil, nil, nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}
	if stderr, err = c.StderrPipe(); err != nil {
		_ = stdout.Close()
		_ = stdin.Close()
		return nil, nil, nil, fmt.Errorf("while opening stderr pipe: %w", err)
	}
	return stdin, stdout, stderr, nil
}

func sendPlan(p *plan.Plan, w io.WriteCloser) error {
	e := toml.NewEncoder(w)
	terr := e.Encode(p)
	ierr := w.Close()
	if terr != nil {
		return fmt.Errorf("while sending input plan: %w", terr)
	}
	if ierr != nil {
		return fmt.Errorf("while closing input pipe: %w", ierr)
	}
	return nil
}
