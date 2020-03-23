// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"context"
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
	OnCopy(src, dst string)

	// OnCopyFinish lets the observer know when a file copy has finished.
	OnCopyFinish()
}

// SFTPer provides a mockable interface for SFTP.
type SFTPer interface {
	// Create tries to create a file at path, and, if successful, opens a File pointing to it.
	Create(path string) (io.WriteCloser, error)
	// MkdirAll recursively makes the directories mentioned in dir.
	MkdirAll(dir string) error
}

// PutMapping copies the files in the (remote-to-local) map mapping to the SFTP client client.
// It checks ctx for cancellation between operations.
func PutMapping(ctx context.Context, client SFTPer, o CopyObserver, mapping map[string]string) error {
	o.OnCopyStart(len(mapping))
	defer o.OnCopyFinish()

	if err := mkdirs(ctx, client, mappingDirs(mapping)); err != nil {
		return err
	}
	return putFiles(ctx, client, o, mapping)
}

func putFiles(ctx context.Context, client SFTPer, o CopyObserver, mapping map[string]string) error {
	for rpath, lpath := range mapping {
		if err := iohelp.CheckDone(ctx); err != nil {
			return err
		}
		if err := sftpPutFile(client, rpath, lpath); err != nil {
			return err
		}
		o.OnCopy(lpath, rpath)
	}
	return nil
}

func mkdirs(ctx context.Context, client SFTPer, dirs []string) error {
	for _, dir := range dirs {
		if err := iohelp.CheckDone(ctx); err != nil {
			return err
		}
		if err := client.MkdirAll(dir); err != nil {
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

func sftpPutFile(cli SFTPer, rpath, lpath string) error {
	r, err := os.Open(filepath.FromSlash(lpath))
	if err != nil {
		return err
	}
	w, err := cli.Create(rpath)
	if err != nil {
		_ = r.Close()
		return err
	}

	_, cperr := io.Copy(w, r)
	wcerr := w.Close()
	rcerr := r.Close()

	if cperr != nil {
		return cperr
	}
	if wcerr != nil {
		return wcerr
	}
	return rcerr
}
