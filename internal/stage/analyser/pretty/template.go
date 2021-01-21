// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"strings"
	"text/template"

	"github.com/c4-project/c4t/internal/subject/obs"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/plan/analysis"
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

	// ShowMutation is true if mutation testing information should be shown.
	ShowMutation bool
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

{{- with .Obs -}}{{- template "interesting" . -}}{{- end -}}

{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- else }}  No subject outcomes available.
{{- end -}}
`

	tmplMutations = `
{{- range $mut, $analysis := .Analysis.Mutation }}  ## Mutant {{ $mut }}
{{ range $analysis }}    - {{ .HitBy }}: {{ .NumHits }} hit(s){{ if .Killed }} *KILLED*{{ end }}
{{ end -}}
{{- else }}  No mutations were enabled.
{{- end -}}
`

	tmplRoot = `
{{- if .ShowPlanInfo -}}
# Plan
{{ template "planInfo" . }}
{{ end -}}
{{- if .ShowCompilers -}}
# Compilers
{{ template "compilers" . -}}
{{/* 'compilers' adds its own newline */}}{{- end -}}
{{- if .ShowSubjects -}}
# Subject Outcomes
{{ template "byStatus" . -}}
{{/* 'byStatus' adds its own newline */}}{{- end -}}
{{- if .ShowMutation -}}
# Mutation Testing
{{ template "mutations" . -}}
{{/* 'mutations' adds its own newline */}}{{- end -}}
`
)

func indent(n int) string {
	return "        " + strings.Repeat("  ", n+1)
}

func getTemplate() (*template.Template, error) {
	t, err := template.New("root").Parse(tmplRoot)
	if err != nil {
		return nil, err
	}
	if t, err = obs.AddObsTemplates(t, indent); err != nil {
		return nil, err
	}
	return iohelp.ParseTemplateStrings(t, map[string]string{
		"timeset":        tmplTime,
		"byStatus":       tmplByStatus,
		"compilers":      tmplCompilers,
		"compilerCounts": tmplCompilerCounts,
		"compilerInfo":   tmplCompilerInfo,
		"compilerLogs":   tmplCompilerLogs,
		"mutations":      tmplMutations,
		"planInfo":       tmplPlanInfo,
		"stages":         tmplStages,
	})
}
