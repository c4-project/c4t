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

	tmplCompilerCounts = `{{ range $status, $count := . }}  {{ $status }}: {{ $count }}
{{ end -}}`

	tmplTime = `{{ if . }}time sec min={{ .Min.Seconds }} avg={{ .Mean.Seconds }} max={{ .Max.Seconds }}{{ else }}no time report{{ end }}`

	tmplCompilers = `
{{- range $cname, $compiler := .Compilers -}}
compiler {{ $cname }} ({{ .Info }}):
  {{ template "timeset" .Time }}
{{ template "compilerCounts" .Counts }}
{{- end -}}
`

	tmplByStatus = `
{{- range $status, $corpus := .ByStatus -}}
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

{{ template "compilers" . }}
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
		"timeset":        tmplTime,
		"states":         tmplStateset,
		"byStatus":       tmplByStatus,
		"compilers":      tmplCompilers,
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
