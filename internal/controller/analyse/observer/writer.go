// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"fmt"
	"io"
	"text/template"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/plan/analysis"
)

// AnalysisWriter provides the ability to output human-readable summaries of analyses to a writer.
type AnalysisWriter struct {
	w    io.Writer
	tmpl *template.Template
}

// NewAnalysisWriter constructs an analysis writer using writer outw.
func NewAnalysisWriter(outw io.Writer) (*AnalysisWriter, error) {
	t, err := getTemplate()
	if err != nil {
		return nil, err
	}

	return &AnalysisWriter{w: iohelp.EnsureWriter(outw), tmpl: t}, nil
}

// Write writes an unsourced analysis an to this writer.
func (a *AnalysisWriter) Write(an analysis.Analysis) error {
	return a.tmpl.ExecuteTemplate(a.w, "root", an)
}

// OnAnalysis writes an unsourced analysis an to this writer; if an error occurs, it tries to rescue.
func (a *AnalysisWriter) OnAnalysis(an analysis.Analysis) {
	if err := a.Write(an); err != nil {
		a.handleError(err)
	}
}

// OnSave does nothing (for now).
func (a *AnalysisWriter) OnSave(Saving) {
}

// OnSaveFileMissing does nothing (for now).
func (a *AnalysisWriter) OnSaveFileMissing(Saving, string) {
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
