// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/c4-project/c4t/internal/stat"
)

//go:embed template/*.tmpl
var templates embed.FS

func getTemplate() (*template.Template, error) {
	dir, err := fs.Sub(templates, "template")
	if err != nil {
		return nil, err
	}
	return template.ParseFS(dir, "*.tmpl")
}

// context is the root structure visible in the stats pretty-printer.
type context struct {
	Stats        *stat.Set
	MutantFilter stat.MutantFilter
	UseTotals    bool
}

// Span gets from m the span required by the context.
func (c context) Span(m stat.Machine) stat.MachineSpan {
	if c.UseTotals {
		return m.Total
	}
	return m.Session
}
