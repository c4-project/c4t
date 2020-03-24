// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/transfer"
)

// Send translates p to the remote host, using SFTP to copy over its harness files.
func (r *SSHRunner) Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	n := transfer.NewNormaliser(r.runner.Config.DirCopy)
	rp := *p
	var err error
	if rp.Corpus, err = n.Corpus(rp.Corpus); err != nil {
		return nil, err
	}

	// We only send the harnesses, to avoid wasting SFTP bandwidth.
	return &rp, r.sendMapping(ctx, n.MappingsOfKind(transfer.NKHarness))
}

func (r *SSHRunner) sendMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := remote.SendMapping(ctx, (*remote.SFTPCopier)(cli), r.observer, ms)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
