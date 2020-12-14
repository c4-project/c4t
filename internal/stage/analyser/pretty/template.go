// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"text/template"

	"github.com/MattWindsor91/c4t/internal/subject/obs"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"

	"github.com/MattWindsor91/c4t/internal/plan/analysis"
)

// WriteContext is the type of roots sent to the template engine.
type WriteContext struct {
	// The analysis to write.
	Analysis *analysis.Analysis

	// ShowCompilers is true if compiler breakdowns should be shown.
	ShowCompilers bool

	// ShowCompilerLogs is true if compiler logs should be shown.
	ShowCompilerLogs bool

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

	tmplCompilerCounts = `{{ range $status, $count := . }}      - {{ $status }}: {{ $count }} subject(s)
{{ end -}}`

	tmplTime = `{{ if . }}Min {{ .Min.Seconds }} Avg {{ .Mean.Seconds }} Max {{ .Max.Seconds }}{{ else }}N/A{{ end }}`

	tmplCompilerInfo = `    - style: {{ .Style }}
    - arch: {{ .Arch }}
    - opt: {{ if .SelectedOpt -}}{{ if .SelectedOpt.Name }}{{ .SelectedOpt.Name }}{{ else }}none{{ end -}}{{ else }}none{{- end }}
    - mopt: {{ if .SelectedMOpt }}{{ .SelectedMOpt }}{{ else }}none{{ end }}`

	tmplCompilerLogs = `
{{- range $sname, $log := . }}{{ if $log }}      #### {{ $sname }}
` + "```" + `
{{ $log }}
` + "```" + `
{{ end -}}
{{ end }}
`

	tmplCompilers = `
{{- range $cname, $compiler := .Analysis.Compilers }}  ## {{ $cname }}
{{ template "compilerInfo" .Info }}
    ### Times (sec)
      - compile: {{ template "timeset" .Time }}
      - run: {{ template "timeset" .RunTime }}
    ### Results
{{ template "compilerCounts" .Counts }}
{{- if $.ShowCompilerLogs }}    ### Logs
{{ template "compilerLogs" .Logs }}
{{ end -}}
{{ end -}}
`

	tmplByStatus = `
{{- range $status, $corpus := .Analysis.ByStatus -}}
{{- if (and $corpus (or (not $status.IsOk) $.ShowOk)) }}  ## {{ $status }} ({{ len $corpus }})
{{ range $sname, $subject := $corpus }}    - {{ $sname }}
{{ range $compiler, $compile := .Compilations -}}
{{- with .Run -}}
{{- if eq $status .Status }}      - {{ $compiler }}
{{ end -}}

{{- with .Obs -}}{{- template "obs" . -}}{{- end -}}

{{- end -}}
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
	if t, err = obs.AddObsTemplates(t, func(n int) string { return "" }); err != nil {
		return nil, err
	}
	return iohelp.ParseTemplateStrings(t, map[string]string{
		"timeset":        tmplTime,
		"byStatus":       tmplByStatus,
		"compilers":      tmplCompilers,
		"compilerCounts": tmplCompilerCounts,
		"compilerInfo":   tmplCompilerInfo,
		"compilerLogs":   tmplCompilerLogs,
		"planInfo":       tmplPlanInfo,
		"stages":         tmplStages,
	})
}
