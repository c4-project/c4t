// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"text/template"

	"github.com/MattWindsor91/act-tester/internal/plan/analysis"
)

// WriteContext is the type of roots sent to the template engine.
type WriteContext struct {
	// The analysis to write.
	Analysis *analysis.Analysis

	// ShowCompilers is true if compiler breakdowns should be shown.
	ShowCompilers bool

	// ShowOk is true if subjects with the 'ok' status should be shown.
	ShowOk bool

	// ShowPlanInfo is true if plan metadata should be shown.
	ShowPlanInfo bool

	// ShowSubjects is true if subject information should be shown.
	ShowSubjects bool
}

const (
	tmplStages = `  ## Stages
{{- range . }}
    - {{ .Stage }}: completed {{ .CompletedOn }}, took {{ .Duration.Seconds }} sec(s)
{{- end -}}
`

	tmplPlanInfo = `  - created at: {{ .Analysis.Plan.Metadata.Creation }}
  - seed: {{ .Analysis.Plan.Metadata.Seed }}
  - version: {{ .Analysis.Plan.Metadata.Version }}
{{ template "stages" .Analysis.Plan.Metadata.Stages -}}
`

	tmplStateset = `{{ range . }}          {{ range $k, $v := . }}  {{ $k }} = {{ $v }}{{- end }}
{{ end -}}`

	tmplObs = `{{- if .CounterExamples -}}
{{- if .Witnesses }}        - witnessing observations:
{{ template "states" .Witnesses }}
{{ end }}        - counter-example observations:
{{ template "states" .CounterExamples }}
{{- end -}}`

	tmplCompilerCounts = `{{ range $status, $count := . }}      - {{ $status }}: {{ $count }} subject(s)
{{ end -}}`

	tmplTime = `{{ if . }}Min {{ .Min.Seconds }} Avg {{ .Mean.Seconds }} Max {{ .Max.Seconds }}{{ else }}N/A{{ end }}`

	tmplCompilerInfo = `    - style: {{ .Style }}
    - arch: {{ .Arch }}
    - opt: {{ if .SelectedOpt -}}{{ if .SelectedOpt.Name }}{{ .SelectedOpt.Name }}{{ else }}none{{ end -}}{{ else }}none{{- end }}
    - mopt: {{ if .SelectedMOpt }}{{ .SelectedMOpt }}{{ else }}none{{ end }}`

	tmplCompilers = `
{{- range $cname, $compiler := .Analysis.Compilers }}  ## {{ $cname }}
{{ template "compilerInfo" .Info }}
    ### Times (sec)
      - compile: {{ template "timeset" .Time }}
      - run: {{ template "timeset" .RunTime }}
    ### Results
{{ template "compilerCounts" .Counts -}}
{{ end -}}
`

	tmplByStatus = `
{{- range $status, $corpus := .Analysis.ByStatus -}}
{{- if (and $corpus (or (not $status.IsOk) $.ShowOk)) }}  ## {{ $status }} ({{ len $corpus }})
{{ range $sname, $subject := $corpus }}    - {{ $sname }}
{{ range $compiler, $compile := .Runs -}}

{{- if eq $status .Status }}      - {{ $compiler }}
{{ end -}}

{{- if .Obs -}}{{- template "obs" .Obs -}}{{- end -}}

{{- end -}}
{{- end }}
{{ end -}}
{{- else }}  No subject outcomes available.
{{- end -}}
`

	tmplRoot = `
{{- if .ShowPlanInfo -}}
# Plan
{{ template "planInfo" . }}
{{ end -}}
{{- if .ShowCompilers -}}
# Compilers
{{ template "compilers" . }}
{{ end -}}
{{- if .ShowSubjects -}}
# Subject Outcomes
{{ template "byStatus" . }}
{{ end -}}
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
		"planInfo":       tmplPlanInfo,
		"stages":         tmplStages,
	} {
		t, err = t.New(n).Parse(ts)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
