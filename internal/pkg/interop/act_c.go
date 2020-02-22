package interop

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// BinActC is the name of the ACT C services binary.
const BinActC = "act-c"

// ErrSubjectNil is an error returned if one calls ProbeSubject(nil).
var ErrSubjectNil = errors.New("subject pointer is nil")

// ProbeSubject populates s with information gleaned from investigating its litmus file.
func (a *ActRunner) ProbeSubject(s *subject.Subject) error {
	if s == nil {
		return ErrSubjectNil
	}

	var h Header
	if err := a.DumpHeader(&h, s.Litmus); err != nil {
		return fmt.Errorf("header read on %s failed: %w", s.Litmus, err)
	}
	s.Name = h.Name

	var st Statset
	if err := a.DumpStats(&st, s.Litmus); err != nil {
		return fmt.Errorf("stats read on %s failed: %w", s.Litmus, err)
	}
	s.Threads = st.Threads

	return nil
}

// DumpHeader runs act-c dump-header on the subject at path, writing the results to h.
func (a *ActRunner) DumpHeader(h *Header, path string) error {
	var obuf bytes.Buffer
	sargs := StandardArgs{Verbose: false}

	cmd := a.Command(BinActC, "dump-header", sargs, path)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return err
	}

	return h.Read(&obuf)
}

// DumpStats runs act-c dump-stats on the subject at path, writing the stats to s.
func (a *ActRunner) DumpStats(s *Statset, path string) error {
	var obuf bytes.Buffer
	sargs := StandardArgs{Verbose: false}

	cmd := a.Command(BinActC, "dump-stats", sargs, path)
	cmd.Stdout = &obuf
	// TODO(@MattWindsor91): allow redirecting this
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return s.Parse(&obuf)
}
