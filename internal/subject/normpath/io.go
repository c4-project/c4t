// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package normpath

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/1set/gut/yos"
	"github.com/c4-project/c4t/internal/helper/errhelp"
)

// ErrMissingFile occurs when we request a subject file but it isn't available.
var ErrMissingFile = errors.New("subject file not available")

// ReadSubjectFile is a low-level function for reading a subject file with path path, relative to directory root.
//
// If the file is present on disk at root/path, this behaves like os.ReadFile.  Otherwise, if the path starts
// with a directory DIR, it assumes that a file root/DIR.tar.gz exists containing path, and attempts to load the file
// from there.  (This is the convention used by the saver when saving tarballs of subjects.)
//
// path may be a slashpath.
func ReadSubjectFile(root, path string) ([]byte, error) {
	// it seems that filepath.Clean subsumes filepath.FromSlash?
	path = filepath.Clean(path)

	apath := filepath.Join(root, path)
	if !yos.ExistFile(apath) {
		return readSubjectFileFromArchive(root, path)
	}
	return os.ReadFile(apath)
}

func readSubjectFileFromArchive(root, path string) ([]byte, error) {
	tp, ok := tarPath(root, path)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingFile, path)
	}
	// Technically this is redundant, as trying to open the file will error if it doesn't exist;
	// we do this just to get a more accurate 'side-path didn't work' error message.
	// TODO(@MattWindsor91): Merging the 'missing file' and 'open error' cases into a compound error would be nice.
	if !yos.ExistFile(tp) {
		return nil, fmt.Errorf("%w: %s (tried archive %s, but couldn't find it)", ErrMissingFile, path, tp)
	}
	return readSubjectFileFromTar(tp, path)
}

func tarPath(root, path string) (string, bool) {
	// Assuming that the first directory in the path is also the name of the tarball.
	tarName, ok := getFirstDir(path)
	return filepath.Join(root, tarName+TarSuffix), ok
}

func readSubjectFileFromTar(tpath string, path string) ([]byte, error) {
	f, err := os.Open(tpath)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	tr := tar.NewReader(gr)
	out, err := readSubjectFileFromTarReader(tr, tpath, path)
	gcerr := gr.Close()
	fcerr := f.Close()
	return out, errhelp.FirstError(err, gcerr, fcerr)
}

func readSubjectFileFromTarReader(tr *tar.Reader, tarFile, rpath string) ([]byte, error) {
	for {
		hd, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("%w: %s (not present in archive %s)", ErrMissingFile, rpath, tarFile)
		}
		if err != nil {
			return nil, err
		}

		// TODO(@MattWindsor): case insensitivity?
		if filepath.Clean(hd.Name) == rpath {
			return io.ReadAll(tr)
		}
	}
}

func getFirstDir(path string) (dir string, ok bool) {
	// TODO(@MattWindsor91): is this safe?
	pfrags := strings.SplitN(path, string(filepath.Separator), 2)
	if len(pfrags) != 2 {
		return "", false
	}
	return pfrags[0], true
}
