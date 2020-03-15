// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// LocalMach runs the 'mach' binary locally.
type LocalMach struct {
	Dir string
}

// Run runs act-tester-mach locally.
func (m *LocalMach) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	// TODO(@MattWindsor91): remote invocation
	// TODO(@MattWindsor91): observation
	cmd := exec.CommandContext(ctx, "act-tester-mach", "-d", m.Dir)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("while opening stdin pipe: %w", err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		_ = in.Close()
		return nil, fmt.Errorf("while opening stdout pipe: %w", err)
	}

	tin := toml.NewEncoder(in)

	if err := cmd.Start(); err != nil {
		_ = out.Close()
		_ = in.Close()
		return nil, fmt.Errorf("while starting local runner: %w", err)
	}

	terr := tin.Encode(p)
	ierr := in.Close()
	if terr != nil {
		_ = cmd.Wait()
		return nil, fmt.Errorf("while sending input plan: %w", terr)
	}
	if ierr != nil {
		_ = cmd.Wait()
		return nil, fmt.Errorf("while closing input pipe: %w", ierr)
	}

	var p2 plan.Plan
	if _, err := toml.DecodeReader(out, &p2); err != nil {
		_ = cmd.Wait()
		return nil, fmt.Errorf("while decoding the output plan: %w", err)
	}

	// Waiting _should_ close the pipes.
	werr := cmd.Wait()
	return &p2, werr
}
