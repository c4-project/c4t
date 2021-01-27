// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs

import (
	"io"
	"strings"
	"text/template"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

const (
	// tmplStateset is a template that prints out a set of observation states.
	//
	// It requires the funcmap to have a function 'obsIndent' from an integer indent level to an indent string.
	tmplStateset = `
{{- range . -}}
{{- obsIndent 1 -}}
{{- if .Occurrences -}}[{{ .Occurrences }}x] {{ end -}}
{{- with $sv := .Values -}}
{{- range $j, $v := .Vars -}}
  {{ if ne $j 0 }}, {{ end }}{{ $v }} = {{ index $sv $v }}
{{- else -}}
{{- end -}}
{{- end }}{{/* deliberate newline here */}}
{{ end -}}`

	// tmplObsInteresting is a template that outputs 'interesting' states on an observation.
	//
	// It requires the funcmap to have a function 'obsIndent' from an integer indent level to an indent string.
	tmplObsInteresting = `
{{- if .Flags.IsPartial -}}
{{ obsIndent 0 }}WARNING: this observation is partial.
{{ end -}}

{{- if .Flags.IsInteresting -}}
{{- if .Flags.IsExistential -}}

{{ obsIndent 0 }}postcondition witnessed by
{{- with .Witnesses -}}:
{{ template "stateset" . }}
{{- else }} at least one of these states:
{{ template "stateset" .States }}
{{- end -}}

{{- else -}}

{{ obsIndent 0 }}postcondition violated by
{{- with .CounterExamples }}:
{{ template "stateset" . }}
{{- else }} at least one of these states:
{{ template "stateset" .States }}
{{- end -}}

{{- end -}}
{{- end -}}
`

	tmplDnf = `forall (
{{- range $i, $s := .States }}
  {{ if eq $i 0 }}  {{ else }}\/{{ end }} (
{{- with $vs := .Values -}}
{{- range $j, $v := .Vars -}}
  {{ if ne $j 0 }} /\ {{ end }}{{ $v }} == {{ index $vs $v }}
{{- else -}}
  true
{{- end -}}
{{- end -}}
)
{{- else }}
  true
{{- end }}
)
`

	tmplPretty = `
{{- if .ShowInteresting -}}{{ template "interesting" .Obs }}{{- end -}}
{{- if and .ShowInteresting .Mode.Dnf }}{{/* TODO(@MattWindsor91): make this unnecessary */}}
postcondition covering all observed states:

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

func (p prettyContext) ShowInteresting() bool {
	return p.Mode.Interesting && (p.Obs.Flags.IsPartial() || p.Obs.Flags.IsInteresting())
}

// AddObsTemplates adds to t a set of templates useful for pretty-printing observations.
//
// indent is a function that should indent observation lines n places, as well as adding any indenting needed to put the
// lines into context.
func AddObsTemplates(t *template.Template, indent func(n int) string) (*template.Template, error) {
	t = t.Funcs(template.FuncMap{"obsIndent": indent})
	return iohelp.ParseTemplateStrings(t, map[string]string{
		"stateset":    tmplStateset,
		"interesting": tmplObsInteresting,
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
