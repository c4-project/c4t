// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package query

import "text/template"

const (
	tmplStateset = `{{ range . }}      {{ range $k, $v := . }}  {{ $k }} = {{ $v }}{{- end }}
{{ end -}}`

	tmplObs = `{{- if .CounterExamples }}      valid final states observed:
{{ template "states" .Witnesses }}
      counter-examples observed:
{{ template "states" .CounterExamples }}
  {{- end -}}`

	tmplCompilerCounts = `
{{- $compilers := .Compilers -}}
{{- range $cname, $compiler := .Analysis.CompilerCounts -}}
compiler {{ $cname }} ({{ index $compilers $cname }}):
{{ range $status, $count := $compiler }}  {{ $status }}: {{ $count }}
{{ end -}}
{{- end -}}
`

	tmplByStatus = `
{{- range $status, $corpus := .Analysis.ByStatus -}}
{{- if not $status.IsOk -}}
{{- if $corpus -}}
status {{ $status }}:
{{ range $sname, $subject := $corpus }}  subject {{ $sname }}:
{{ range $compiler, $compile := .Runs -}}

{{- if eq $status .Status }}    {{ $compiler }}
{{ end -}}

{{- if .Obs -}}{{- template "obs" .Obs -}}{{- end -}}

{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

	tmplRoot = `COMPILER BREAKDOWN:

{{ template "compilerCounts" . }}
SUBJECT REPORT:

{{ template "byStatus" . -}}
`
)

func getTemplate() (*template.Template, error) {
	t, err := template.New("root").Parse(tmplRoot)
	if err != nil {
		return nil, err
	}
	for n, ts := range map[string]string{
		"states":         tmplStateset,
		"byStatus":       tmplByStatus,
		"compilerCounts": tmplCompilerCounts,
		"obs":            tmplObs,
	} {
		t, err = t.New(n).Parse(ts)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
