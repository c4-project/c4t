// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"io"
	"strings"
	"text/template"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
)

const (
	// TmplStateset is a template that prints out a set of observation states.
	//
	// It requires the funcmap to have a function 'obsIndent' from an integer indent level to an indent string.
	TmplStateset = `
{{- range $s := . -}}
{{- obsIndent 1 -}}
{{- range $j, $v := .Vars -}}
  {{ if ne $j 0 }}, {{ end }}{{ $v }} = {{ index $s $v }}
{{- else -}}
{{- end }}{{/* deliberate newline here */}}
{{ end -}}`

	// TmplObsInteresting is a template that outputs 'interesting' states on an observation.
	//
	// It requires the funcmap to have a function 'obsIndent' from an integer indent level to an indent string.
	TmplObsInteresting = `
{{- if .Flags.IsInteresting -}}
{{- if .Flags.IsExistential -}}

{{- with .Witnesses }}{{ obsIndent 0}}- existential witnessed by:
{{ template "stateset" . }}
{{ end -}}

{{- else -}}

{{- with .CounterExamples }}{{ obsIndent 0 }}- postcondition violated by:
{{ template "stateset" . }}
{{ end -}}

{{- end -}}
{{- end -}}
`

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
{{- if .Mode.Interesting -}}{{ template "interesting" .Obs }}{{- end -}}
{{- if and .Mode.Interesting .Mode.Dnf }}{{/* TODO(@MattWindsor91): make this unnecessary */}}
{{ end }}
{{- if .Mode.Dnf -}}{{ template "dnf" .Obs }}{{ end -}}
`
)

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

func AddObsTemplates(t *template.Template, indent func(n int) string) (*template.Template, error) {
	t = t.Funcs(template.FuncMap{"obsIndent": indent})
	return iohelp.ParseTemplateStrings(t, map[string]string{
		"stateset":    TmplStateset,
		"interesting": TmplObsInteresting,
	})
}

func makeTemplate() (*template.Template, error) {
	t, err := template.New("root").Parse(tmplPretty)
	if err != nil {
		return nil, err
	}
	if t, err = AddObsTemplates(t, func(n int) string { return strings.Repeat("  ", n) }); err != nil {
		return nil, err
	}
	return iohelp.ParseTemplateStrings(t,
		map[string]string{
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
