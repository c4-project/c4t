// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compilation

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/1set/gut/yos"
)

// TODO(MattWindsor91): a lot of this probably belongs elsewhere, but import cycles and methods make it difficult.

// TarSuffix is the extension used by the saver when saving tarballs, and presumed by archive-transparent file readers.
const TarSuffix = ".tar.gz"

var (
	// ErrNoCompilerLog occurs when we ask for the compiler log of a subject that doesn't have one.
	ErrNoCompilerLog = errors.New("compiler result has no log file")
	// ErrMissingFile occurs when we request a subject file but it isn't available.
	ErrMissingFile = errors.New("subject file not available")
)

// ReadLog tries to read in the log for compiler, taking paths relative to root.
// If the compiler log doesn't exist relative to root, and its path is of the form FOO/BAR, we assume that it is in a
// saved tarball called FOO.tar.gz (as file BAR) in root, and attempt to extract it.
func (c *CompileFileset) ReadLog(root string) ([]byte, error) {
	if ystring.IsBlank(c.Log) {
		return nil, ErrNoCompilerLog
	}
	return readSubjectFile(root, c.Log)
}

func readSubjectFile(root, path string) ([]byte, error) {
	// it seems that filepath.Clean subsumes filepath.FromSlash?
	path = filepath.Clean(path)

	apath := filepath.Join(root, path)
	if !yos.ExistFile(apath) {
		return readSubjectFileFromArchive(root, path)
	}
	return ioutil.ReadFile(apath)
}

func readSubjectFileFromArchive(root, path string) ([]byte, error) {
	sname, rpath, ok := splitSubjectPrefix(path)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingFile, path)
	}
	tarPath := filepath.Join(root, sname+TarSuffix)
	// Technically this is redundant, as trying to open the file will error if it doesn't exist;
	// we do this just to get a more accurate 'side-path didn't work' error message.
	// TODO(@MattWindsor91): Merging the 'missing file' and 'open error' cases into a compound error would be nice.
	if !yos.ExistFile(tarPath) {
		return nil, fmt.Errorf("%w: %s (tried archive %s, but couldn't find it)", ErrMissingFile, path, tarPath)
	}
	return readSubjectFileFromTar(tarPath, rpath)
}

func readSubjectFileFromTar(tarPath string, rpath string) ([]byte, error) {
	f, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	tr := tar.NewReader(gr)
	out, err := readSubjectFileFromTarReader(tr, tarPath, rpath)
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
		if filepath.Clean(hd.Name) != rpath {
			continue
		}
		return ioutil.ReadAll(tr)
	}
}

func splitSubjectPrefix(path string) (prefix, rest string, ok bool) {
	// TODO(@MattWindsor91): is this safe?
	pfrags := strings.SplitN(path, string(filepath.Separator), 2)
	if len(pfrags) != 2 {
		return "", "", false
	}
	return pfrags[0], pfrags[1], true
}
