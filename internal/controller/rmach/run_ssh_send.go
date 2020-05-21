// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Send translates p to the remote host, using SFTP to copy over its harness files.
func (r *SSHRunner) Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	n := normaliser.NewCorpus(r.runner.Config.DirCopy)
	rp := *p
	var err error
	if rp.Corpus, err = n.Normalise(rp.Corpus); err != nil {
		return nil, err
	}

	// We only send the harness source code, to avoid wasting SFTP bandwidth.
	return &rp, r.sendMapping(ctx, n.Mappings.RenamesMatching(filekind.C, filekind.InHarness))
}

func (r *SSHRunner) sendMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := remote.SendMapping(ctx, (*remote.SFTPCopier)(cli), ms, r.observers...)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
