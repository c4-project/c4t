// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"fmt"
	"os"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Pathset contains the paths for a lifter.
type Pathset struct {
	// TODO(@MattWindsor91): can anything else be done here?
	root string

	paths map[string]map[string]string
}

// Pather abstracts over the path resolution for a lifter.
type Pather interface {
	// Prepare sets up a pathset to deal with the architecture IDs arches and subject names subjects.
	Prepare(arches []id.ID, subjects []string) error

	// Path gets the path to the directory prepared for arch and subject.
	// It fails if no such directory has been prepared.
	Path(arch id.ID, subject string) (string, error)
}

//go:generate mockery --name=Pather

// NewPathset makes a pathset under root.
func NewPathset(root string) *Pathset {
	return &Pathset{root: root, paths: nil}
}

// Prepare sets up a pathset to deal with the architecture IDs arches and subject names subjects.
func (p *Pathset) Prepare(arches []id.ID, subjects []string) error {
	p.paths = make(map[string]map[string]string, len(arches))
	for _, a := range arches {
		as := a.String()
		p.paths[as] = make(map[string]string, len(subjects))
		for _, s := range subjects {
			segs := p.pathSegs(a, s)
			p.paths[as][s] = path.Join(segs...)

			if err := os.MkdirAll(p.paths[as][s], 0744); err != nil {
				return err
			}
		}
	}
	return nil
}

// Path gets the path to the directory prepared for arch and subject.
// It fails if no such directory has been prepared.
func (p *Pathset) Path(arch id.ID, subject string) (string, error) {
	as := arch.String()
	amap, ok := p.paths[as]
	if !ok {
		return "", fmt.Errorf("arch %s not prepared", as)
	}
	dir, ok := amap[subject]
	if !ok {
		return "", fmt.Errorf("subject %s not prepared for arch %s", subject, as)
	}
	return dir, nil
}

func (p *Pathset) pathSegs(a id.ID, s string) []string {
	segs := make([]string, 1, 2+len(a.Tags()))
	segs[0] = p.root
	segs = append(append(segs, a.Tags()...), s)
	return segs
}
