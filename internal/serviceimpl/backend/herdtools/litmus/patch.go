// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/helper/errhelp"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/1set/gut/yos"
)

const (
	// IncludeStdbool is the include directive that the litmus patcher will insert if needed.
	// It is exported for testing purposes.
	IncludeStdbool = "#include <stdbool.h>"

	atomicCast = "(_Atomic int)"

	dumpPrefix = "fprintf(fhist,"
)

// patch patches the Litmus files in p, which originated from a Litmus invocation in inFile.
func (l *Instance) patch() error {
	if !l.Fixset.NeedsPatch() {
		return nil
	}

	rpath := l.mainCFile()
	wpath, err := l.patchToTemp(rpath)
	if err != nil {
		return err
	}

	// NOT os.Rename, which doesn't work between filesystems.
	return yos.MoveFile(wpath, rpath)
}

// mainCFile guesses where Litmus is going to put the main file in its generated harness.
func (l *Instance) mainCFile() string {
	file := iohelp.ExtlessFile(l.Job.In.Litmus.Path) + ".c"
	return filepath.Join(l.Job.Out.Dir, file)
}

func (l *Instance) patchToTemp(rpath string) (wpath string, err error) {
	r, rerr := os.Open(rpath)
	if rerr != nil {
		return "", fmt.Errorf("can't open C file for reading: %w", rerr)
	}
	wpath, werr := l.patchReaderToTemp(r)
	cerr := r.Close()
	return wpath, errhelp.FirstError(werr, cerr)
}

func (l *Instance) patchReaderToTemp(r io.Reader) (string, error) {
	w, werr := ioutil.TempFile("", "*.c")
	if werr != nil {
		return "", fmt.Errorf("can't open temp file for reading: %w", werr)
	}
	wpath := w.Name()
	// Right now, there's only one thing to patch, so this is fairly easy.
	err := l.Fixset.PatchMainFile(r, w)
	cerr := w.Close()
	return wpath, errhelp.FirstError(err, cerr)
}

// PatchMainFile patches the main C file represented by rw according to this fixset.
func (f *Fixset) PatchMainFile(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		if err := f.patchLine(w, sc.Text()); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (f *Fixset) patchLine(w io.Writer, line string) error {
	line = f.patchWithinLine(line)

	if _, err := fmt.Fprintln(w, line); err != nil {
		return fmt.Errorf("can't write to patched file: %w", err)
	}

	if strings.Contains(line, "/* Includes */") {
		if err := f.patchStdbool(w); err != nil {
			return fmt.Errorf("can't insert include into buffer: %w", err)
		}
	}
	return nil
}

// patchWithinLine does line-level patches on line.
func (f *Fixset) patchWithinLine(line string) string {
	switch {
	case f.RemoveAtomicCasts && isDump(line):
		return strings.ReplaceAll(line, atomicCast, "")
	default:
		return line
	}
}

// isDump is a heuristic for checking whether line is the problematic dumping fprintf in the Litmus harness that
// includes atomic casts.
func isDump(line string) bool {
	ls := strings.TrimSpace(line)
	return strings.HasPrefix(ls, dumpPrefix)
}

// patchStdbool inserts an include for stdbool into w.
func (f *Fixset) patchStdbool(w io.Writer) error {
	if !f.InjectStdbool {
		return nil
	}

	_, err := io.WriteString(w, IncludeStdbool+"\n")
	return err
}
