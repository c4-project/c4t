// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package query

import "text/template"

const (
	tmplStateset = `{{ range . }}  {{ range $k, $v := . }}  {{ $k }} = {{ $v }}{{- end }}
{{ end -}}`

	tmplObs = `{{- if .CounterExamples }}  valid final states observed:
{{ template "states" .Witnesses }}
  counter-examples observed:
{{ template "states" .CounterExamples }}
  {{- end -}}`

	tmplCorpus = `
{{- $compilers := .Compilers -}}
{{- range $sname, $subject := .Corpus -}}
  {{- range $compiler, $compile := .Runs -}}

{{- if not .Status.IsOk -}}
{{ $sname }} {{ index $compilers $compiler }}: {{ .Status }}
{{ end -}}

{{- if .Obs -}}{{- template "obs" .Obs -}}{{- end -}}

  {{- end -}}
{{- end -}}
`

	tmplPlan = `
{{- template "corpus" . -}}
`
)

func getTemplate() (*template.Template, error) {
	t, err := template.New("plan").Parse(tmplPlan)
	if err != nil {
		return nil, err
	}
	for n, ts := range map[string]string{
		"states": tmplStateset,
		"corpus": tmplCorpus,
		"obs":    tmplObs,
	} {
		t, err = t.New(n).Parse(ts)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
