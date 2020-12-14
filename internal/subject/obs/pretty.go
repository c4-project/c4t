// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"io"
	"text/template"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
)

const (
	tmplDnf = `
forall (
{{- range $i, $s := .States }}
  {{ if eq $i 0 }}  {{ else }}\/{{ end }} (
{{- range $j, $v := .Vars -}}
  {{ if ne $j 0 }} /\ {{ end }}{{ $v }} == {{ index $s $v }}
{{- else -}}
  true
{{- end -}}
)
{{- else }}
  true
{{- end }}
)
`

	tmplPretty = `
{{- if .Mode.Dnf -}}{{ template "dnf" .Obs }}{{ end -}}
`
)

// PrettyMode controls various pieces of pretty-printer functionality.
type PrettyMode struct {
	// Dnf controls whether the pretty-printer prints a disjunctive-normal-form postcondition.
	Dnf bool
}

type prettyContext struct {
	Mode PrettyMode
	Obs  Obs
}

func makeTemplate() (*template.Template, error) {
	return iohelp.TemplateFromStrings(tmplPretty, map[string]string{
		"dnf": tmplDnf,
	})
}

// Pretty pretty-prints an observation o onto w according to mode m.
func Pretty(w io.Writer, o Obs, m PrettyMode) error {
	t, err := makeTemplate()
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "root", prettyContext{Mode: m, Obs: o})
}
