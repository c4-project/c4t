// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"errors"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// ErrCollision occurs if the normaliser tries to map two files to the same normalised path.
// Usually, this is an internal error.
var ErrCollision = errors.New("path already mapped by normaliser")

const (
	FileBin        = "a.out"
	FileCompileLog = "compile.log"
	FileOrigLitmus = "orig.litmus"
	FileFuzzLitmus = "fuzz.litmus"
	FileFuzzTrace  = "fuzz.trace"
	DirCompiles    = "compiles"
	DirHarnesses   = "harnesses"
)

// Normaliser contains state necessary to normalise a single subject's paths.
// This is useful for archiving the subject inside a tarball, or copying it to another host.
type Normaliser struct {
	// root is the prefix to add to every normalised name.
	root string

	// Mappings contains maps from normalised names to original names.
	// (The mappings are this way around to help us notice collisions.)
	Mappings map[string]string
}

// NewNormaliser constructs a new Normaliser relative to root.
func NewNormaliser(root string) *Normaliser {
	return &Normaliser{
		root:     root,
		Mappings: make(map[string]string),
	}
}

// Subject normalises mappings from subject component files to 'normalised' names, relative to root.
func (n *Normaliser) Subject(s subject.Subject) (*subject.Subject, error) {
	var err error
	s.Litmus, err = n.replaceAndAdd(s.Litmus, FileOrigLitmus)
	if s.Fuzz != nil && err == nil {
		s.Fuzz, err = n.fuzz(*s.Fuzz)
	}
	if s.Compiles != nil && err == nil {
		s.Compiles, err = n.compiles(s.Compiles)
	}
	if s.Harnesses != nil && err == nil {
		s.Harnesses, err = n.harnesses(s.Harnesses)
	}
	// No need to normalise runs
	return &s, err
}

func (n *Normaliser) fuzz(f subject.Fuzz) (*subject.Fuzz, error) {
	var err error
	if f.Files.Litmus, err = n.replaceAndAdd(f.Files.Litmus, FileFuzzLitmus); err != nil {
		return nil, err
	}
	f.Files.Trace, err = n.replaceAndAdd(f.Files.Trace, FileFuzzTrace)
	return &f, err
}

func (n *Normaliser) harnesses(hs map[string]subject.Harness) (map[string]subject.Harness, error) {
	nhs := make(map[string]subject.Harness, len(hs))
	for archstr, h := range hs {
		var err error
		nhs[archstr], err = n.harness(archstr, h)
		if err != nil {
			return nil, err
		}
	}
	return nhs, nil
}

func (n *Normaliser) harness(archstr string, h subject.Harness) (subject.Harness, error) {
	oldPaths := h.Paths()
	h.Dir = path.Join(n.root, DirHarnesses, archstr)
	for i, np := range h.Paths() {
		if err := n.add(oldPaths[i], np); err != nil {
			return h, err
		}
	}
	return h, nil
}

func (n *Normaliser) compiles(cs map[string]subject.CompileResult) (map[string]subject.CompileResult, error) {
	ncs := make(map[string]subject.CompileResult, len(cs))
	for cidstr, c := range cs {
		var err error
		ncs[cidstr], err = n.compile(cidstr, c)
		if err != nil {
			return nil, err
		}
	}

	return ncs, nil
}

func (n *Normaliser) compile(cidstr string, c subject.CompileResult) (subject.CompileResult, error) {
	var err error
	if c.Files.Bin, err = n.replaceAndAdd(c.Files.Bin, DirCompiles, cidstr, FileBin); err != nil {
		return c, err
	}
	c.Files.Log, err = n.replaceAndAdd(c.Files.Log, DirCompiles, cidstr, FileCompileLog)
	return c, err
}

// replaceAndAdd adds the path assembled by joining segs together as a mapping from opath.
// If opath is empty, this just returns ("", nil) and does no addition.
func (n *Normaliser) replaceAndAdd(opath string, segs ...string) (string, error) {
	if opath == "" {
		return "", nil
	}

	npath := path.Join(segs...)
	err := n.add(opath, npath)
	return npath, err
}

// add tries to add the mapping between opath and npath to the normaliser's mappings.
// It fails if there is a collision.
func (n *Normaliser) add(opath, npath string) error {
	if _, ok := n.Mappings[npath]; ok {
		return fmt.Errorf("%w: %q", ErrCollision, npath)
	}
	n.Mappings[npath] = opath
	return nil
}
