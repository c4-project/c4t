// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"io"
	"os"
	"text/template"

	"github.com/c4-project/c4t/internal/stat"
)

// Printer provides the ability to output human-readable summaries of the statistics file to a writer.
type Printer struct {
	w    io.Writer
	tmpl *template.Template
	ctx  context
}

func NewPrinter(o ...Option) (*Printer, error) {
	t, err := getTemplate()
	if err != nil {
		return nil, err
	}

	aw := &Printer{w: os.Stdout, tmpl: t}
	Options(o...)(aw)

	return aw, nil
}

// Write writes a summary of statistic set s to this writer.
func (p *Printer) Write(s stat.Set) error {
	c := p.ctx
	c.Stats = &s
	return p.tmpl.ExecuteTemplate(p.w, "root.tmpl", c)
}
