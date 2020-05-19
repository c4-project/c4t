// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import "text/template"

const (
	tmplStateset = `{{ range . }}      {{ range $k, $v := . }}  {{ $k }} = {{ $v }}{{- end }}
{{ end -}}`

	tmplObs = `{{- if .CounterExamples -}}
{{- if .Witnesses }}      obs valid
{{ template "states" .Witnesses }}
{{ end }}      obs counter-examples
{{ template "states" .CounterExamples }}
{{- end -}}`

	tmplCompilerCounts = `{{ range $status, $count := . }}  count {{ $status }} {{ $count }}
{{ end -}}`

	tmplTime = `{{ if . }}sec min {{ .Min.Seconds }} avg {{ .Mean.Seconds }} max {{ .Max.Seconds }}{{ else }}n/a{{ end }}`

	tmplCompilerInfo = `style {{ .Style }} arch {{ .Arch }} opt {{ if .SelectedOpt -}}
{{if .SelectedOpt.Name}}{{ .SelectedOpt.Name }}{{ else }}none{{ end -}}{{ else }}none
{{- end }} mopt {{ if .SelectedMOpt }}{{ .SelectedMOpt }}{{ else }}none{{ end }}`

	tmplCompilers = `
{{- range $cname, $compiler := .Compilers -}}
compiler {{ $cname }}
  spec {{ template "compilerInfo" .Info }}
  timings compile {{ template "timeset" .Time }}
  timings run {{ template "timeset" .RunTime }}
{{ template "compilerCounts" .Counts -}}
{{ end -}}
`

	tmplByStatus = `
{{- range $status, $corpus := .ByStatus -}}
{{- if not $status.IsOk -}}
{{- if $corpus -}}
status {{ $status }}
{{ range $sname, $subject := $corpus }}  subject {{ $sname }}
{{ range $compiler, $compile := .Runs -}}

{{- if eq $status .Status }}    compiler {{ $compiler }}
{{ end -}}

{{- if .Obs -}}{{- template "obs" .Obs -}}{{- end -}}

{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
`

	tmplRoot = `COMPILER BREAKDOWN

{{ template "compilers" . }}
SUBJECT REPORT

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
		"compilerInfo":   tmplCompilerInfo,
		"obs":            tmplObs,
	} {
		t, err = t.New(n).Parse(ts)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
