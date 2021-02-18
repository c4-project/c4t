// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"embed"
	"io"
	"io/fs"
	"strings"
	"text/template"
)

//go:embed template
var templates embed.FS

// PrettyMode controls various pieces of pretty-printer functionality.
type PrettyMode struct {
	// Dnf controls whether the pretty-printer prints a disjunctive-normal-form postcondition.
	Dnf bool
	// Interesting controls whether the pretty-printer prints 'interesting' state results.
	Interesting bool
}

type prettyContext struct {
	Mode PrettyMode
	Obs  Obs
}

// ShowSummary gets whether the pretty printer has been configured to show a summary of its current observation.
func (p prettyContext) ShowSummary() bool {
	return p.Mode.Interesting && (p.Obs.Flags.IsPartial() || p.Obs.Flags.IsInteresting())
}

// AddCommonTemplates adds to t a set of templates useful for pretty-printing observations.
//
// indent is a function that should indent observation lines n places, as well as adding any indenting needed to put the
// lines into context.
func AddCommonTemplates(t *template.Template, indent func(n int) string) (*template.Template, error) {
	t = t.Funcs(template.FuncMap{"obsIndent": indent})
	efs, err := fs.Sub(templates, "template/common")
	if err != nil {
		return nil, err
	}
	return t.ParseFS(efs, "*.tmpl")
}

func makeTemplate() (*template.Template, error) {
	t := template.New("root.tmpl")
	efs, err := fs.Sub(templates, "template")
	if err != nil {
		return nil, err
	}
	if t, err = AddCommonTemplates(t, func(n int) string { return strings.Repeat("  ", n) }); err != nil {
		return nil, err
	}
	return t.ParseFS(efs, "*.tmpl")
}

// Pretty pretty-prints an observation o onto w according to mode m.
func Pretty(w io.Writer, o Obs, m PrettyMode) error {
	t, err := makeTemplate()
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "root.tmpl", prettyContext{Mode: m, Obs: o})
}
