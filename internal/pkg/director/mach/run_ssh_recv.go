// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"context"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
	"github.com/MattWindsor91/act-tester/internal/pkg/transfer"
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
	norm := transfer.NewNormaliser(path.Join(r.recvRoot, ls.Name))
	rs, ok := rcorp[ls.Name]
	if !ok {
		return fmt.Errorf("subject not in remote corpus: %s", ls.Name)
	}
	ns, err := norm.Subject(rs)
	if err != nil {
		return fmt.Errorf("can't normalise corpus: %w", err)
	}
	ls.Runs = ns.Runs
	ls.Compiles = ns.Compiles
	return r.recvMapping(ctx, norm.MappingsOfKind(transfer.NKCompile))
}

func (r *SSHRunner) recvMapping(ctx context.Context, ms map[string]string) error {
	cli, err := r.runner.NewSFTP()
	if err != nil {
		return err
	}

	perr := remote.RecvMapping(ctx, (*remote.SFTPCopier)(cli), r.observer, ms)
	cerr := cli.Close()

	if perr != nil {
		return perr
	}
	return cerr
}
