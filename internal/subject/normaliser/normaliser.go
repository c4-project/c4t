// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package normaliser provides utilities for archiving and transferring plans, corpora, and subjects.
package normaliser

import (
	"errors"
	"fmt"
	"path"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/subject/normpath"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/model/recipe"

	"github.com/c4-project/c4t/internal/model/filekind"

	"github.com/c4-project/c4t/internal/subject"
)

// ErrCollision occurs if the normaliser tries to map two files to the same normalised path.
// Usually, this is an internal error.
var ErrCollision = errors.New("path already mapped by normaliser")

// Normaliser contains state necessary to normalise a single subject's paths.
// This is useful for archiving the subject inside a tarball, or copying it to another host.
type Normaliser struct {
	// root is the prefix to add to every normalised name.
	root string

	// err is the first error this normaliser encountered.
	err error

	// Mappings contains maps from normalised names to original names.
	// (The mappings are this way around to help us notice collisions.)
	Mappings Map
}

// New constructs a new Normaliser relative to root.
func New(root string) *Normaliser {
	return &Normaliser{
		root:     root,
		Mappings: make(map[string]Entry),
	}
}

// Normalise normalises mappings from subject component files to 'normalised' names.
func (n *Normaliser) Normalise(s subject.Subject) (*subject.Subject, error) {
	n.err = nil

	s.Source.Path = n.replaceAndAdd(s.Source.Path, filekind.Litmus, filekind.InOrig, normpath.FileOrigLitmus)
	s.Fuzz = n.fuzz(s.Fuzz)
	s.Compilations = n.compilations(s.Compilations)
	s.Recipes = n.recipes(s.Recipes)
	// No need to normalise runs
	return &s, n.err
}

func (n *Normaliser) fuzz(of *subject.Fuzz) *subject.Fuzz {
	if of == nil {
		return nil
	}
	f := *of
	f.Litmus.Path = n.replaceAndAdd(f.Litmus.Path, filekind.Litmus, filekind.InFuzz, normpath.FileFuzzLitmus)
	f.Trace = n.replaceAndAdd(f.Trace, filekind.Trace, filekind.InFuzz, normpath.FileFuzzTrace)
	return &f
}

func (n *Normaliser) recipes(rs recipe.Map) recipe.Map {
	if rs == nil {
		return nil
	}

	nrs := make(recipe.Map, len(rs))
	for arch, r := range rs {
		nrs[arch] = n.recipe(arch, r)
	}
	return nrs
}

func (n *Normaliser) recipe(arch id.ID, h recipe.Recipe) recipe.Recipe {
	oldPaths := h.Paths()
	h.Dir = normpath.RecipeDir(n.root, arch.String())
	for i, np := range h.Paths() {
		n.add(oldPaths[i], np, filekind.GuessFromFile(np), filekind.InRecipe)
	}
	return h
}

func (n *Normaliser) compilations(cs compilation.Map) compilation.Map {
	if cs == nil {
		return nil
	}
	ncs := make(compilation.Map, len(cs))
	for cid, c := range cs {
		ncs[cid] = n.compilation(cid, c)
	}
	return ncs
}

func (n *Normaliser) compilation(compiler id.ID, c compilation.Compilation) compilation.Compilation {
	if c.Compile != nil {
		c.Compile = n.compile(compiler, *c.Compile)
	}
	return c
}

func (n *Normaliser) compile(compiler id.ID, c compilation.CompileResult) *compilation.CompileResult {
	c.Files.Bin = n.replaceAndAdd(c.Files.Bin, filekind.Bin, filekind.InCompile, normpath.DirCompiles, compiler.String(), normpath.FileBin)
	c.Files.Log = n.replaceAndAdd(c.Files.Log, filekind.Log, filekind.InCompile, normpath.DirCompiles, compiler.String(), normpath.FileCompileLog)
	return &c
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
	n.Mappings[npath] = Entry{
		Original: opath,
		Kind:     k,
		Loc:      l,
	}
	return npath
}
