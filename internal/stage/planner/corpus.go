// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner

import (
	"context"

	"github.com/MattWindsor91/c4t/internal/quantity"

	"github.com/1set/gut/yos"

	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"

	"github.com/MattWindsor91/c4t/internal/subject/corpus"

	"github.com/MattWindsor91/c4t/internal/subject"
)

// SubjectProber is the interface of types that allow filling in of subject information.
type SubjectProber interface {
	// ProbeSubject probes the Litmus test at filepath litmus, producing a named subject record.
	ProbeSubject(ctx context.Context, litmus string) (*subject.Named, error)
}

func (p *Planner) planCorpus(ctx context.Context, fs ...string) (corpus.Corpus, error) {
	files, err := ExpandLitmusInputs(fs)
	if err != nil {
		return nil, err
	}
	c := CorpusPlanner{
		Files:      files,
		Prober:     p.source.SProbe,
		Observers:  lowerToBuilder(p.observers),
		Quantities: p.quantities,
	}
	return c.Plan(ctx)
}

// CorpusPlanner contains the state required to plan the corpus part of an initial plan file.
type CorpusPlanner struct {
	// Files contains the files that are to be included in the plan.
	Files []string
	// Observers observe the process of building the corpus.
	Observers []builder.Observer
	// Prober tells the planner how to probe corpus files for specific information.
	Prober SubjectProber
	// Quantities contains the target size and worker count of the corpus.
	Quantities quantity.PlanSet
}

// Plan probes each subject in this planner's corpus file list, producing a Corpus proper.
// It does not sample; sampling is left to the perturb stage.
func (p *CorpusPlanner) Plan(ctx context.Context) (corpus.Corpus, error) {
	cfg := p.makeBuilderConfig()
	// TODO(@MattWindsor91): perform corpus pruning
	return builder.ParBuild(ctx, p.Quantities.NWorkers, corpus.New(p.Files...), cfg,
		func(ctx context.Context, named subject.Named, requests chan<- builder.Request) error {
			// TODO(@MattWindsor91): make it so we don't get the litmus file through the *name* of the subject!
			// TODO(@MattWindsor91): overwrite 'named' with gleaned information
			return p.probeSubject(ctx, named.Name, requests)
		},
	)
}

func (p *CorpusPlanner) makeBuilderConfig() builder.Config {
	return builder.Config{
		Init:      nil,
		Observers: p.Observers,
		Manifest: builder.Manifest{
			Name:  "plan",
			NReqs: len(p.Files),
		},
	}
}

// ExpandLitmusInputs expands a list of possible Litmus files or containing directories into a flat file list.
func ExpandLitmusInputs(in []string) ([]string, error) {
	// TODO(@MattWindsor91): abstract this properly, it's only exposed because the coverage thing uses it too.

	files := make([]string, 0, len(in))
	var err error
	for _, f := range in {
		if files, err = expandFile(f, files); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func expandFile(f string, curr []string) ([]string, error) {
	if yos.ExistDir(f) {
		return expandDir(f, curr)
	}
	// Delegate handling of non-files to the point where we actually open them.
	return append(curr, f), nil
}

func expandDir(d string, curr []string) ([]string, error) {
	ents, err := yos.ListMatch(d, yos.ListIncludeFile, "*.litmus")
	if err != nil {
		return nil, err
	}
	for _, ent := range ents {
		curr = append(curr, ent.Path)
	}
	return curr, nil
}

func (p *CorpusPlanner) probeSubject(ctx context.Context, f string, ch chan<- builder.Request) error {
	s, err := p.Prober.ProbeSubject(ctx, f)
	if err != nil {
		return err
	}
	return builder.AddRequest(s).SendTo(ctx, ch)
}
