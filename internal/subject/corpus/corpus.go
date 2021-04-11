// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package corpus concerns test corpora (collections of named subjects).
package corpus

import (
	"errors"
	"fmt"
	"sort"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/helper/stringhelp"

	"github.com/c4-project/c4t/internal/subject"
)

var (
	// ErrCorpusDup occurs when we try to Add a subject into a corpus under a name that is already taken.
	ErrCorpusDup = errors.New("duplicate corpus entry")

	// ErrMapRename occurs when we try to change the name of an entry inside a Map.
	ErrMapRename = errors.New("tried to rename a corpus entry")

	// ErrSmall occurs when the viable test corpus is smaller than that requested by the user.
	ErrSmall = errors.New("test corpus too small")

	// ErrNone is a variant of ErrSmall that occurs when the viable test corpus is empty.
	ErrNone = fmt.Errorf("%w: no corpus given", ErrSmall)
)

// Corpus is the type of test corpora (groups of test subjects).
type Corpus map[string]subject.Subject

// New creates a blank Corpus from a list of names.
func New(names ...string) Corpus {
	corpus := make(Corpus, len(names))
	for _, n := range names {
		corpus[n] = subject.Subject{}
	}
	return corpus
}

// FilterToNames filters c to contain only subjects whose names are contained within names.
func (c Corpus) FilterToNames(names ...string) Corpus {
	c2 := make(Corpus, len(names))
	nset := stringhelp.NewSet(names...)
	for n, s := range c {
		if _, ok := nset[n]; !ok {
			continue
		}
		c2[n] = s
	}
	return c2
}

// Add tries to add s to the corpus.
// It fails if the corpus already has a subject with the given name.
func (c Corpus) Add(s subject.Named) error {
	_, exists := c[s.Name]
	if exists {
		return fmt.Errorf("%w: %s", ErrCorpusDup, s.Name)
	}

	c[s.Name] = s.Subject
	return nil
}

// Copy makes a deep copy of this corpus.
func (c Corpus) Copy() Corpus {
	cc := make(Corpus, len(c))
	for n, s := range c {
		cc[n] = s
	}
	return cc
}

// Names returns a sorted list of this corpus's subject names.
func (c Corpus) Names() []string {
	ns := make([]string, len(c))
	i := 0
	for n := range c {
		ns[i] = n
		i++
	}
	sort.Strings(ns)
	return ns
}

// EraseCompilations deletes the compilation entries for each subject in c.
func (c Corpus) EraseCompilations() {
	for n, s := range c {
		s.Compilations = compilation.Map{}
		c[n] = s
	}
}
