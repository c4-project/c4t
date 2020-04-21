// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package query

import (
	"context"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
	"github.com/MattWindsor91/act-tester/internal/view"

	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	c "github.com/urfave/cli/v2"
)

func App(outw, errw io.Writer) *c.App {
	return &c.App{
		Name:  "act-tester-query",
		Usage: "performs human-readable queries on a plan file",
		Flags: flags(),
		Action: func(ctx *c.Context) error {
			return run(ctx, outw, errw)
		},
		Writer:                 outw,
		ErrWriter:              errw,
		HideHelpCommand:        true,
		UseShortOptionHandling: true,
	}
}

func flags() []c.Flag {
	return []c.Flag{
		stdflag.PlanFileCliFlag(),
		// TODO(@MattWindsor91): template stuff
	}
}

type query struct {
	outw io.Writer
}

const tmpl = `
{{- define "states" -}}
{{ range . }}  {{ range $k, $v := . }}  {{ $k }} = {{ $v }}{{- end }}
{{ end -}}
{{- end -}}

{{- $compilers := .Compilers -}}
{{- range $sname, $subject := .Corpus -}}
  {{- range $compiler, $compile := .Runs -}}

{{- if not .Status.IsOk -}}
{{ $sname }} {{ index $compilers $compiler }}: {{ .Status }}
{{ end -}}

  {{- if and .Obs .Obs.CounterExamples }}  valid final states observed:
{{ template "states" .Obs.Witnesses }}
  counter-examples observed:
{{ template "states" .Obs.CounterExamples }}
  {{- end -}}

  {{- end -}}
{{- end -}}
`

func (q query) Run(_ context.Context, p *plan.Plan) (*plan.Plan, error) {
	t, err := template.New("plan").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	err = t.Execute(q.outw, p)

	return p, err
}

func run(ctx *c.Context, outw io.Writer, _ io.Writer) error {
	pf := stdflag.PlanFileFromCli(ctx)
	q := query{outw: outw}
	return view.RunOnPlanFile(ctx.Context, q, pf, ioutil.Discard)
}
