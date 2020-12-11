// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"io"
	"text/template"
)

const (
	prettyTmpl = `
{{- if .Mode.Dnf -}}
forall (
{{- range $i, $s := .Obs.States }}
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
{{ end -}}
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

// Pretty pretty-prints an observation o onto w according to mode m.
func Pretty(w io.Writer, o Obs, m PrettyMode) error {
	t, err := template.New("obs").Parse(prettyTmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, prettyContext{Mode: m, Obs: o})
}
