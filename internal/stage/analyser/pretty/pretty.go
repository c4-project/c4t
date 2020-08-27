// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package pretty provides a pretty-printer for analyses.
package pretty

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
)

// Printer provides the ability to output human-readable summaries of analyses to a writer.
type Printer struct {
	w    io.Writer
	tmpl *template.Template
	ctx  WriteContext
}

// NewPrinter constructs a pretty-printer using options o.
func NewPrinter(o ...Option) (*Printer, error) {
	t, err := getTemplate()
	if err != nil {
		return nil, err
	}

	aw := &Printer{w: os.Stdout, tmpl: t}
	Options(o...)(aw)

	return aw, nil
}

// Write writes an unsourced analysis a to this writer.
func (p *Printer) Write(a analysis.Analysis) error {
	c := p.ctx
	c.Analysis = &a
	return p.tmpl.ExecuteTemplate(p.w, "root", c)
}

// OnAnalysis writes an unsourced analysis a to this printer; if an error occurs, it tries to rescue.
func (p *Printer) OnAnalysis(a analysis.Analysis) {
	if err := p.Write(a); err != nil {
		p.handleError(err)
	}
}

// WriteSourced writes a sourced analysis a to this printer.
func (p *Printer) WriteSourced(a analysis.WithRun) error {
	if _, err := fmt.Fprintf(p.w, "# %s #\n\n", &a.Run); err != nil {
		return err
	}
	return p.Write(a.Analysis)
}

func (p *Printer) handleError(err error) {
	_, _ = fmt.Fprintf(p.w, "ERROR OUTPUTTING ANALYSIS: %s\n", err)
}
