// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

import (
	"context"
	"fmt"

	copy2 "github.com/MattWindsor91/act-tester/internal/copier"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"github.com/MattWindsor91/act-tester/internal/plan"
)

// Recv copies bits of remp into locp, including run information and any compiler failures.
// It uses SFTP to transfer back any compile logs.
func (r *RemoteRunner) Recv(ctx context.Context, locp, remp *plan.Plan) (*plan.Plan, error) {
	locp.Metadata.Stages = remp.Metadata.Stages

	norm := normaliser.NewCorpus(r.localRoot)
	ncorp, err := norm.Normalise(remp.Corpus)
	if err != nil {
		return nil, fmt.Errorf("can't normalise corpus: %w", err)
	}

	if err := r.mergeSubjects(locp, ncorp); err != nil {
		return nil, err
	}
	return locp, r.recvMapping(ctx, norm.Mappings.RenamesMatching(filekind.Any, filekind.InCompile))
}

func (r *RemoteRunner) mergeSubjects(locp *plan.Plan, rcorp corpus.Corpus) error {
	return locp.Corpus.Map(func(sn *subject.Named) error {
		return r.mergeSubject(sn, rcorp)
	})
}

func (r *RemoteRunner) mergeSubject(ls *subject.Named, rcorp corpus.Corpus) error {
	rs, ok := rcorp[ls.Name]
	if !ok {
		return fmt.Errorf("subject not in remote corpus: %s", ls.Name)
	}
	ls.Runs = rs.Runs
	ls.Compiles = rs.Compiles
	return nil
}

func (r *RemoteRunner) recvMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := copy2.RecvMapping(ctx, (*copy2.SFTP)(cli), ms, r.observers...)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
