// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/subject/normaliser"
)

// Send translates p to the remote host, using SFTP to copy over any recipe files.
func (r *RemoteRunner) Send(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	n := normaliser.NewCorpus(r.runner.Config.DirCopy)
	rp := *p
	var err error
	if rp.Corpus, err = n.Normalise(rp.Corpus); err != nil {
		return nil, err
	}

	// We only send the recipe source code, to avoid wasting SFTP bandwidth.
	// TODO(@MattWindsor91): actually check which files are mentioned in recipe instructions?
	return &rp, r.sendMapping(ctx, n.Mappings.RenamesMatching(filekind.C, filekind.InRecipe))
}

func (r *RemoteRunner) sendMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := copy2.SendMapping(ctx, (*copy2.SFTP)(cli), ms, r.observers...)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
