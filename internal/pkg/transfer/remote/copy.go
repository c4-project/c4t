// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/pkg/helpers/iohelp"
)

// CopyObserver is an interface for types that observe an SFTP file copy.
type CopyObserver interface {
	// OnCopyStart lets the observer know when a file copy (of nfiles files) is beginning.
	OnCopyStart(nfiles int)

	// OnCopy lets the observer know that a file copy (from path src to path dst) has happened.
	OnCopy(dst, src string)

	// OnCopyFinish lets the observer know when a file copy has finished.
	OnCopyFinish()
}

// OnCopyStart sends an OnCopyStart observation to multiple observers.
func OnCopyStart(nfiles int, cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopyStart(nfiles)
	}
}

// OnCopy sends an OnCopy observation to multiple observers.
func OnCopy(dst, src string, cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopy(dst, src)
	}
}

// OnCopyFinish sends an OnCopyFinish observation to multiple observers.
func OnCopyFinish(cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopyFinish()
	}
}

// Copier provides a mockable interface for remote copying.
type Copier interface {
	// Create tries to create a file at path, and, if successful, opens a write-closer pointing to it.
	Create(path string) (io.WriteCloser, error)
	// Open tries to open a file at path, and, if successful, opens a read-closer pointing to it.
	Open(path string) (io.ReadCloser, error)
	// MkdirAll recursively makes the directories mentioned in dir.
	MkdirAll(dir string) error
}

// LocalCopier implements Copier through os.
type LocalCopier struct{}

// Create calls os.Create on path.
func (l LocalCopier) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Open calls os.Open on path.
func (l LocalCopier) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// MkkdirAll calls os.MkdirAll on path, with vaguely sensible permissions.
func (l LocalCopier) MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0744)
}

// SendMapping is shorthand for CopyMapping where the source is a LocalCopier.
func SendMapping(ctx context.Context, dst Copier, mapping map[string]string, o ...CopyObserver) error {
	return CopyMapping(ctx, dst, LocalCopier{}, mapping, o...)
}

// RecvMapping is shorthand for CopyMapping where the source is a LocalCopier.
func RecvMapping(ctx context.Context, src Copier, mapping map[string]string, o ...CopyObserver) error {
	return CopyMapping(ctx, LocalCopier{}, src, mapping, o...)
}

// CopyMapping copies the files in the (dest-to-src) map mapping from dst to src.
// It checks ctx for cancellation between operations.
func CopyMapping(ctx context.Context, dst, src Copier, mapping map[string]string, o ...CopyObserver) error {
	OnCopyStart(len(mapping), o...)
	defer OnCopyFinish(o...)

	if err := mkdirs(ctx, dst, mappingDirs(mapping)); err != nil {
		return err
	}
	return copyFiles(ctx, dst, src, mapping, o...)
}

func copyFiles(ctx context.Context, dst, src Copier, mapping map[string]string, o ...CopyObserver) error {
	for dpath, spath := range mapping {
		if err := iohelp.CheckDone(ctx); err != nil {
			return err
		}
		if err := copyFile(dst, src, dpath, spath); err != nil {
			return fmt.Errorf("copying %s to %s: %w", spath, dpath, err)
		}
		OnCopy(dpath, spath, o...)
	}
	return nil
}

func mkdirs(ctx context.Context, dst Copier, dirs []string) error {
	for _, dir := range dirs {
		if err := iohelp.CheckDone(ctx); err != nil {
			return err
		}
		if err := dst.MkdirAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func mappingDirs(mapping map[string]string) []string {
	dirs := make([]string, len(mapping))
	i := 0
	for k := range mapping {
		dirs[i] = path.Dir(k)
		i++
	}
	return dirs
}

func copyFile(dst, src Copier, rpath, lpath string) error {
	r, err := src.Open(filepath.FromSlash(lpath))
	if err != nil {
		return err
	}
	w, err := dst.Create(rpath)
	if err != nil {
		_ = r.Close()
		return err
	}

	_, err = iohelp.CopyClose(w, r)
	return err
}
