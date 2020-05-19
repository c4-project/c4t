// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse

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

// NewAnalysisWriter constructs an analysis writer using config c.
func NewAnalysisWriter(c *Config) (*AnalysisWriter, error) {
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	t, err := getTemplate()
	if err != nil {
		return nil, err
	}

	return &AnalysisWriter{w: iohelp.EnsureWriter(c.Out), tmpl: t}, nil
}

// Write writes an unsourced analysis an to this writer.
func (a *AnalysisWriter) Write(an *analysis.Analysis) error {
	if an == nil {
		return nil
	}

	return a.tmpl.ExecuteTemplate(a.w, "root", an)
}

// WriteSourced writes a sourced analysis an to this writer.
func (a *AnalysisWriter) WriteSourced(an *analysis.Sourced) error {
	if an == nil {
		return nil
	}

	if _, err := fmt.Fprintf(a.w, "# %s #\n\n", &an.Run); err != nil {
		return err
	}

	return a.Write(an.Collation)
}
