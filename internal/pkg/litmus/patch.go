package litmus

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// IncludeStdbool is the include directive that the litmus patcher will insert if needed.
// It is exported for testing purposes.
const IncludeStdbool = "#include <stdbool.h>"

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
	if err := f.patchStdbool(w); err != nil {
		return fmt.Errorf("can't insert include into buffer: %w", err)
	}

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("can't copy C file: %w", err)
	}

	return nil
}

func (f *Fixset) patchStdbool(w io.Writer) error {
	if !f.InjectStdbool {
		return nil
	}

	_, err := io.WriteString(w, IncludeStdbool+"\n")
	return err
}
