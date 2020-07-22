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

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// AnalysisWriter provides the ability to output human-readable summaries of analyses to a writer.
type AnalysisWriter struct {
	w    io.Writer
	tmpl *template.Template
	ctx  WriteContext
}

// NewAnalysisWriter constructs an analysis writer using options o.
func NewAnalysisWriter(o ...Option) (*AnalysisWriter, error) {
	t, err := getTemplate()
	if err != nil {
		return nil, err
	}

	aw := &AnalysisWriter{w: os.Stdout, tmpl: t}
	Options(o...)(aw)

	return aw, nil
}

// Write writes an unsourced analysis an to this writer.
func (a *AnalysisWriter) Write(an analysis.Analysis) error {
	c := a.ctx
	c.Analysis = &an
	return a.tmpl.ExecuteTemplate(a.w, "root", c)
}

// OnAnalysis writes an unsourced analysis an to this writer; if an error occurs, it tries to rescue.
func (a *AnalysisWriter) OnAnalysis(an analysis.Analysis) {
	if err := a.Write(an); err != nil {
		a.handleError(err)
	}
}

// OnArchive does nothing (for now).
func (a *AnalysisWriter) OnArchive(observer.ArchiveMessage) {
}

// WriteSourced writes a sourced analysis an to this writer.
func (a *AnalysisWriter) WriteSourced(an analysis.Sourced) error {
	if _, err := fmt.Fprintf(a.w, "# %s #\n\n", &an.Run); err != nil {
		return err
	}
	return a.Write(an.Analysis)
}

func (a *AnalysisWriter) handleError(err error) {
	_, _ = fmt.Fprintf(a.w, "ERROR OUTPUTTING ANALYSIS: %s\n", err)
}
