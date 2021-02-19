// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"embed"
	"io/fs"
	"strings"
	"text/template"
	"time"

	"github.com/c4-project/c4t/internal/subject/obs"
)

//go:embed template
var templates embed.FS

// Config is the type of pretty-printer configuration.
type Config struct {
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

// WithConfig is the type of things wrapped with pretty-printer config.
type WithConfig struct {
	// Data is the wrapped data.
	Data interface{}
	// Config is the configuration.
	Config Config
}

// AddConfig wraps the item x with the config c.
func AddConfig(x interface{}, c Config) WithConfig {
	return WithConfig{Data: x, Config: c}
}

func indent(n int) string {
	return "        " + strings.Repeat("  ", n+1)
}

func getTemplate() (*template.Template, error) {
	t := template.New("root.tmpl")
	efs, err := fs.Sub(templates, "template")
	if err != nil {
		return nil, err
	}
	if t, err = obs.AddCommonTemplates(t, indent); err != nil {
		return nil, err
	}
	return t.Funcs(template.FuncMap{
		"withConfig": AddConfig,
		"time":       func(t time.Time) string { return t.Format(time.StampMilli) },
	}).ParseFS(efs, "*.tmpl")
}
