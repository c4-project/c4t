// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// Read reads plan information from r into p.
func Read(r io.Reader, p *Plan) error {
	_, err := toml.DecodeReader(r, p)
	return err
}

// ReadFile reads plan information from the file named by path into p.
func ReadFile(path string, p *Plan) error {
	_, err := toml.DecodeFile(path, p)
	return err
}

// Write dumps plan p to w.
func (p *Plan) Write(w io.Writer) error {
	enc := toml.NewEncoder(w)
	enc.Indent = "  "
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
