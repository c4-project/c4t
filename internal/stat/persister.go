// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/copier"

	"github.com/c4-project/c4t/internal/machine"

	"github.com/c4-project/c4t/internal/helper/errhelp"
	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/director"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/stage/analyser/saver"
)

// Persister is a forward handler that maintains and persists a statistics set on disk.
type Persister struct {
	// set is the statistics set being persisted.
	set Statset
	// f is the target file (we need a file to be able to truncate properly).
	f *os.File
	// enc is the encoder writing to f.
	enc *json.Encoder
	// err holds any error that has occurred while persisting statistics.
	err error
	// lastCount is the value of EventCount when set was last committed to disk.
	lastCount uint64
}

// NewPersister creates a Persister that reads and writes statistics from f.
// If f is non-empty, it immediately tries to read any existing stats dump, and fails if this doesn't work.
// The Persister takes ownership of f; close f with the Persister's Close method.
func NewPersister(f *os.File) (*Persister, error) {
	sp := Persister{f: f, enc: json.NewEncoder(f)}
	// Not strictly necessary, but makes eyeballing the stats easier.
	sp.enc.SetIndent("", "\t")
	if err := sp.tryReadStats(); err != nil {
		return nil, fmt.Errorf("while reloading stats: %w", err)
	}
	sp.set.ResetForSession()
	return &sp, nil
}

// Close closes this Persister, returning any errors arising from either the stats persisting or file close.
func (s *Persister) Close() error {
	perr := s.err
	cerr := s.f.Close()
	return errhelp.FirstError(perr, cerr)
}

// OnMachine feeds the information from m into the stats set.
func (s *Persister) OnMachines(m machine.Message) {
	s.set.OnMachines(m)
	s.flush()
}

// OnMachine feeds the information from m into the stats set.
func (s *Persister) OnPrepare(m director.PrepareMessage) {
	s.set.OnPrepare(m)
	s.flush()
}

// OnCycle feeds the information from c into the stats set.
func (s *Persister) OnCycle(c director.CycleMessage) {
	s.set.OnCycle(c)
	s.flush()
}

// OnCycleInstance feeds the information from c and m into the stats set.
func (s *Persister) OnCycleInstance(c director.Cycle, m director.InstanceMessage) {
	s.set.OnCycleInstance(c, m)
	s.flush()
}

// OnCycleAnalysis feeds the information from a into the stats set.
func (s *Persister) OnCycleAnalysis(a director.CycleAnalysis) {
	s.set.OnCycleAnalysis(a)
	s.flush()
}

// OnCycleBuild feeds the information from c and m into the stats set.
func (s *Persister) OnCycleBuild(c director.Cycle, m builder.Message) {
	s.set.OnCycleBuild(c, m)
	s.flush()
}

// OnCycleCompiler feeds the information from c and m into the stats set.
func (s *Persister) OnCycleCompiler(c director.Cycle, m compiler.Message) {
	s.set.OnCycleCompiler(c, m)
	s.flush()
}

// OnCycleCopy feeds the information from c and m into the stats set.
func (s *Persister) OnCycleCopy(c director.Cycle, m copier.Message) {
	s.set.OnCycleCopy(c, m)
	s.flush()
}

// OnCycleSave feeds the information from c and m into the stats set.
func (s *Persister) OnCycleSave(c director.Cycle, m saver.ArchiveMessage) {
	s.set.OnCycleSave(c, m)
	s.flush()
}

func (s *Persister) tryReadStats() error {
	if empty, err := iohelp.IsFileEmpty(s.f); err != nil {
		return fmt.Errorf("while checking file for existing stats: %w", err)
	} else if empty {
		s.set.Init()
		return nil
	}
	return s.set.Load(s.f)
}

func (s *Persister) flush() {
	if s.err != nil {
		return
	}
	// Possible ABA problem, but should be ok as long as at most one event happens before each flush.
	if s.set.EventCount == s.lastCount {
		return
	}
	s.lastCount = s.set.EventCount
	if _, s.err = s.f.Seek(0, io.SeekStart); s.err != nil {
		return
	}
	if s.err = s.f.Truncate(0); s.err != nil {
		return
	}
	s.err = s.enc.Encode(s.set)
}

// OpenStatFile opens a file in the appropriate mode for using it as a statistics target.
func OpenStatFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
}
