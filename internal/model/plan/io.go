// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// Ext is the file extension, if any, that should be used for plan files.
const Ext = ".json"

// Read reads plan information from r into p.
func Read(r io.Reader, p *Plan) error {
	return json.NewDecoder(r).Decode(p)
}

// ReadFile reads plan information from the file named by path into p.
func ReadFile(path string, p *Plan) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening plan file %q: %w", path, err)
	}
	perr := Read(r, p)
	cerr := r.Close()
	return iohelp.FirstError(perr, cerr)
}

// Write dumps plan p to w.
func (p *Plan) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(p)
}

// WriteFile dumps plan p to the file named by path.
func (p *Plan) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating plan file: %w", err)
	}
	err = p.Write(f)
	cerr := f.Close()
	return iohelp.FirstError(err, cerr)
}
