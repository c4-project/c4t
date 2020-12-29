// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package copy provides a mockable interface for network-transparent(ish) file copying, and implementations thereof.
package copier

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// Copier provides a mockable interface for remote copying.
type Copier interface {
	// Create tries to create a file at path, and, if successful, opens a write-closer pointing to it.
	Create(path string) (io.WriteCloser, error)
	// Open tries to open a file at path, and, if successful, opens a read-closer pointing to it.
	Open(path string) (io.ReadCloser, error)
	// MkdirAll recursively makes the directories mentioned in dir.
	MkdirAll(dir string) error
}

//go:generate mockery --name=Copier

// SendMapping is shorthand for CopyMapping where the source is a Local.
func SendMapping(ctx context.Context, dst Copier, mapping map[string]string, o ...Observer) error {
	return CopyMapping(ctx, dst, Local{}, mapping, o...)
}

// RecvMapping is shorthand for CopyMapping where the source is a Local.
func RecvMapping(ctx context.Context, src Copier, mapping map[string]string, o ...Observer) error {
	return CopyMapping(ctx, Local{}, src, mapping, o...)
}

// CopyMapping copies the files in the (dest-to-src) map mapping from dst to src.
// It checks ctx for cancellation between operations.
func CopyMapping(ctx context.Context, dst, src Copier, mapping map[string]string, o ...Observer) error {
	OnCopyStart(len(mapping), o...)
	defer OnCopyEnd(o...)

	if err := mkdirs(ctx, dst, mappingDirs(mapping)); err != nil {
		return err
	}
	return copyFiles(ctx, dst, src, mapping, o...)
}

func copyFiles(ctx context.Context, dst, src Copier, mapping map[string]string, o ...Observer) error {
	i := 0
	for dpath, spath := range mapping {
		if err := iohelp.CheckDone(ctx); err != nil {
			return err
		}
		if err := copyFile(dst, src, dpath, spath); err != nil {
			return fmt.Errorf("copying %s to %s: %w", spath, dpath, err)
		}
		OnCopyStep(i, dpath, spath, o...)
		i++
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
