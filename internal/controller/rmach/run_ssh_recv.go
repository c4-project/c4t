// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/remote"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// Recv copies bits of remp into locp, including run information and any compiler failures.
// It uses SFTP to transfer back any compile logs.
func (r *SSHRunner) Recv(ctx context.Context, locp, remp *plan.Plan) (*plan.Plan, error) {
	err := locp.Corpus.Map(func(sn *subject.Named) error {
		return r.recvSubject(ctx, sn, remp.Corpus)
	})
	return locp, err
}

func (r *SSHRunner) recvSubject(ctx context.Context, ls *subject.Named, rcorp corpus.Corpus) error {
	norm := normaliser.New(path.Join(r.recvRoot, ls.Name))
	rs, ok := rcorp[ls.Name]
	if !ok {
		return fmt.Errorf("subject not in remote corpus: %s", ls.Name)
	}
	ns, err := norm.Normalise(rs)
	if err != nil {
		return fmt.Errorf("can't normalise subject: %w", err)
	}
	ls.Runs = ns.Runs
	ls.Compiles = ns.Compiles
	return r.recvMapping(ctx, norm.Mappings.RenamesMatching(filekind.Any, filekind.InCompile))
}

func (r *SSHRunner) recvMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := remote.RecvMapping(ctx, (*remote.SFTPCopier)(cli), ms, r.observers...)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
