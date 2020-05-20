// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package normalise provides utilities for archiving and transferring plans, corpora, and subjects.
package normalise

import (
	"errors"
	"fmt"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
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
	// DirCompiles is the normalised directory for compile results.
	DirCompiles = "compiles"
	// DirHarnesses is the normalised directory for harness results.
	DirHarnesses = "harnesses"
)

// Normaliser contains state necessary to normalise a single subject's paths.
// This is useful for archiving the subject inside a tarball, or copying it to another host.
type Normaliser struct {
	// root is the prefix to add to every normalised name.
	root string

	// err is the first error this normaliser encountered.
	err error

	// Mappings contains maps from normalised names to original names.
	// (The mappings are this way around to help us notice collisions.)
	Mappings map[string]Normalisation
}

// Normalisation is a record in the normaliser's mappings.
// This exists mainly to make it possible to use a Normaliser to work out how to copy a plan to another host,
// but only copy selective subsets of files.
type Normalisation struct {
	// Original is the original path.
	Original string
	// Kind is the kind of path to which this mapping belongs.
	Kind filekind.Kind
	// Loc is an abstraction of the location of the path to which this mapping belongs.
	Loc filekind.Loc
}

// NewNormaliser constructs a new Normaliser relative to root.
func NewNormaliser(root string) *Normaliser {
	return &Normaliser{
		root:     root,
		Mappings: make(map[string]Normalisation),
	}
}

// MappingsMatching filters this normaliser's map to only the files matching kind k and location l.
func (n *Normaliser) MappingsMatching(k filekind.Kind, l filekind.Loc) map[string]string {
	fs := make(map[string]string)
	for n, m := range n.Mappings {
		if m.Kind.Matches(k) && m.Loc.Matches(l) {
			fs[n] = m.Original
		}
	}
	return fs
}

// Corpus normalises mappings for each subject in c.
func (n *Normaliser) Corpus(c corpus.Corpus) (corpus.Corpus, error) {
	c2 := make(corpus.Corpus, len(c))
	for name, s := range c {
		// The aliasing of Mappings here is deliberate.
		snorm := Normaliser{root: path.Join(n.root, name), Mappings: n.Mappings}
		ns, err := snorm.Subject(s)
		if err != nil {
			return nil, fmt.Errorf("normalising %s: %w", name, err)
		}
		c2[name] = *ns
	}
	return c2, nil
}

// Subject normalises mappings from subject component files to 'normalised' names.
func (n *Normaliser) Subject(s subject.Subject) (*subject.Subject, error) {
	s.OrigLitmus = n.replaceAndAdd(s.OrigLitmus, filekind.Litmus, filekind.InOrig, FileOrigLitmus)
	if s.Fuzz != nil {
		s.Fuzz = n.fuzz(*s.Fuzz)
	}
	if s.Compiles != nil {
		s.Compiles = n.compiles(s.Compiles)
	}
	if s.Harnesses != nil {
		s.Harnesses = n.harnesses(s.Harnesses)
	}
	// No need to normalise runs
	return &s, n.err
}

func (n *Normaliser) fuzz(f subject.Fuzz) *subject.Fuzz {
	f.Files.Litmus = n.replaceAndAdd(f.Files.Litmus, filekind.Litmus, filekind.InFuzz, FileFuzzLitmus)
	f.Files.Trace = n.replaceAndAdd(f.Files.Trace, filekind.Trace, filekind.InFuzz, FileFuzzTrace)
	return &f
}

func (n *Normaliser) harnesses(hs map[string]subject.Harness) map[string]subject.Harness {
	nhs := make(map[string]subject.Harness, len(hs))
	for archstr, h := range hs {
		nhs[archstr] = n.harness(archstr, h)
	}
	return nhs
}

func (n *Normaliser) harness(archstr string, h subject.Harness) subject.Harness {
	oldPaths := h.Paths()
	h.Dir = path.Join(n.root, DirHarnesses, archstr)
	for i, np := range h.Paths() {
		n.add(oldPaths[i], np, filekind.GuessFromFile(np), filekind.InHarness)
	}
	return h
}

func (n *Normaliser) compiles(cs map[string]subject.CompileResult) map[string]subject.CompileResult {
	ncs := make(map[string]subject.CompileResult, len(cs))
	for cidstr, c := range cs {
		ncs[cidstr] = n.compile(cidstr, c)
	}
	return ncs
}

func (n *Normaliser) compile(cidstr string, c subject.CompileResult) subject.CompileResult {
	c.Files.Bin = n.replaceAndAdd(c.Files.Bin, filekind.Bin, filekind.InCompile, DirCompiles, cidstr, FileBin)
	c.Files.Log = n.replaceAndAdd(c.Files.Log, filekind.Log, filekind.InCompile, DirCompiles, cidstr, FileCompileLog)
	return c
}

// replaceAndAdd adds the path assembled by joining segs together as a mapping from opath.
// If opath is empty, this just returns "" and does no addition.
func (n *Normaliser) replaceAndAdd(opath string, k filekind.Kind, l filekind.Loc, segs ...string) string {
	if n.err != nil || opath == "" {
		return ""
	}
	return n.add(opath, path.Join(n.root, path.Join(segs...)), k, l)
}

// add tries to add the mapping between opath and npath to the normaliser's mappings, returning npath.
// It fails if there is a collision.
func (n *Normaliser) add(opath, npath string, k filekind.Kind, l filekind.Loc) string {
	if _, ok := n.Mappings[npath]; ok {
		n.err = fmt.Errorf("%w: %q", ErrCollision, npath)
		return npath
	}
	n.Mappings[npath] = Normalisation{
		Original: opath,
		Kind:     k,
		Loc:      l,
	}
	return npath
}
