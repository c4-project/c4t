// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/collate"
	"github.com/MattWindsor91/act-tester/internal/pkg/director/pathset"
	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// Save contains the state used when saving 'interesting' subjects.
type Save struct {
	// Logger is the logger to use when announcing that we're saving subjects.
	Logger *log.Logger

	// NWorkers is the number of workers to use for the collator.
	NWorkers int

	// Paths contains the pathset used to save subjects for a particular machine.
	Paths *pathset.Saved
}

// Run runs the saving stage over plan p.
// It returns p unchanged; this is for signature compatibility with the other director stages.
func (s *Save) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	if s.Paths == nil {
		return nil, iohelp.ErrPathsetNil
	}
	s.Logger = iohelp.EnsureLog(s.Logger)

	if err := s.Paths.Prepare(); err != nil {
		return nil, err
	}

	coll, err := collate.Collate(ctx, p.Corpus, s.NWorkers)
	if err != nil {
		return nil, fmt.Errorf("when collating: %w", err)
	}

	// TODO(@MattWindsor91): find a better way of doing this - this is quite verbose...
	cases := []struct {
		corp corpus.Corpus
		tf   func(string, time.Time) string
	}{
		{corp: coll.CompileFailures, tf: s.Paths.CompileFailureTarFile},
		{corp: coll.RunFailures, tf: s.Paths.RunFailureTarFile},
		{corp: coll.Timeouts, tf: s.Paths.TimeoutTarFile},
		{corp: coll.Flagged, tf: s.Paths.FlaggedTarFile},
	}
	for _, c := range cases {
		if err := s.tarSubjects(c.corp, c.tf, p.Header.Creation); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *Save) tarSubjects(corp corpus.Corpus, tf func(string, time.Time) string, creation time.Time) error {
	for name, sub := range corp {
		tarpath := tf(name, creation)
		if err := os.MkdirAll(filepath.Dir(tarpath), 0744); err != nil {
			return nil
		}
		s.Logger.Printf("archiving %s (to %q)", name, tarpath)
		if err := s.tarSubject(sub, tarpath); err != nil {
			return fmt.Errorf("tarring subject %s: %w", name, err)
		}
	}
	return nil
}

func (s *Save) tarSubject(sub subject.Subject, tarpath string) error {
	tarfile, err := os.Create(tarpath)
	if err != nil {
		return fmt.Errorf("create %s: %w", tarpath, err)
	}
	gzw := gzip.NewWriter(tarfile)
	tarw := tar.NewWriter(gzw)

	werr := s.tarSubjectToWriter(sub, tarw)
	terr := tarw.Close()
	gerr := gzw.Close()

	if werr != nil {
		return werr
	}
	if terr != nil {
		return fmt.Errorf("closing tar: %w", terr)
	}
	if gerr != nil {
		return fmt.Errorf("closing gzip: %w", gerr)
	}
	return nil
}

func (s *Save) tarSubjectToWriter(sub subject.Subject, tarw *tar.Writer) error {
	fs, err := filesToTar(sub)
	if err != nil {
		return err
	}
	for wpath, rpath := range fs {
		if err := s.tarFileToWriter(rpath, wpath, tarw); err != nil {
			return fmt.Errorf("archiving %q: %w", rpath, err)
		}
	}
	return nil
}

func filesToTar(s subject.Subject) (map[string]string, error) {
	n := NewNormaliser("")
	if _, err := n.Subject(s); err != nil {
		return nil, err
	}
	return n.Mappings, nil
}

// tarFileToWriter tars the file at rpath to wpath within the tar archive represented by tarw.
// If rpath is empty, no tarring occurs.
func (s *Save) tarFileToWriter(rpath, wpath string, tarw *tar.Writer) error {
	if rpath == "" {
		return nil
	}

	hdr, err := tarFileHeader(rpath, wpath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.Logger.Println("file missing when archiving error:", rpath)
			return nil
		}
		return err
	}
	if err := tarw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}
	f, err := os.Open(rpath)
	if err != nil {
		return fmt.Errorf("opening %s: %w", rpath, err)
	}
	_, err = io.Copy(tarw, f)
	cerr := f.Close()
	if err != nil {
		return fmt.Errorf("archiving %s: %w", rpath, err)
	}
	return cerr
}

func tarFileHeader(rpath, wpath string) (*tar.Header, error) {
	info, err := os.Stat(rpath)
	if err != nil {
		return nil, fmt.Errorf("can't stat %s: %w", rpath, err)
	}
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return nil, fmt.Errorf("can't get header for %s: %w", rpath, err)
	}
	hdr.Name = wpath
	return hdr, nil
}
