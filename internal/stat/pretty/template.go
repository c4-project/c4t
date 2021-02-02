// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"text/template"

	"github.com/c4-project/c4t/internal/stat"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// TODO(@MattWindsor91): embed these (and any other templates) when Go 1.16 releases.
const (
	tmplMutant = `
{{- .Kills }} kill(s) on {{ .Selections }} selection(s) [{{ .Hits }} hits]
{{- range $s, $k := .Statuses -}}, {{ $s }}={{ $k }}{{ end -}}
`

	tmplMachines = `
{{- with $ctx := . -}}
{{ range $mid, $mach := .Stats.Machines }}  ## {{ $mid }}
{{ with $muts := ($ctx.Span $mach).Mutation -}}
{{ if $ctx.MutantFilter }}    ### Mutants
{{ range $mut := .MutantsWhere $ctx.MutantFilter }}      {{ $mut }}. {{ template "mutant" (index $muts.ByMutant $mut) }}
{{ else }}      No mutants available matching filter.
{{ end -}}
{{- end -}}
{{ else }}    No records available for this machine.
{{- end -}}
{{ else }}  No machines available.
{{ end -}}
{{- end -}}
`

	tmplRoot = `# Machine Report
{{ template "machines" . -}}
`
)

func getTemplate() (*template.Template, error) {
	t, err := template.New("root").Parse(tmplRoot)
	if err != nil {
		return nil, err
	}
	return iohelp.ParseTemplateStrings(t, map[string]string{
		"machines": tmplMachines,
		"mutant":   tmplMutant,
	})
}

// context is the root structure visible in the stats pretty-printer.
type context struct {
	Stats        *stat.Set
	MutantFilter stat.MutantFilter
	UseTotals    bool
}

// Span gets from m the span required by the context.
func (c context) Span(m stat.Machine) stat.MachineSpan {
	if c.UseTotals {
		return m.Total
	}
	return m.Session
}
