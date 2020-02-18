package litmus

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// IncludeStdbool is the include directive that the litmus patcher will insert if needed.
	// It is exported for testing purposes.
	IncludeStdbool = "#include <stdbool.h>"

	atomicCast = "(_Atomic int)"

	dumpPrefix = "fprintf(fhist,"
)

// patch patches the Litmus files in p, which originated from a Litmus invocation in inFile.
func (l *Litmus) patch() error {
	// Right now, there's only one thing to patch, so this is fairly easy.

	if !l.Fixset.InjectStdbool {
		return nil
	}

	path := l.Pathset.MainCFile()

	r, rerr := os.Open(path)
	if rerr != nil {
		return fmt.Errorf("can't open C file for reading: %w", rerr)
	}
	w, werr := ioutil.TempFile("", "*.c")
	if werr != nil {
		_ = r.Close()
		return fmt.Errorf("can't open temp file for reading: %w", rerr)
	}
	wpath := w.Name()
	if err := l.Fixset.PatchMainFile(r, w); err != nil {
		_ = r.Close()
		_ = w.Close()
		return err
	}
	if err := r.Close(); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		_ = w.Close()
		return err
	}
	return os.Rename(wpath, path)
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
