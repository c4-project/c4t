// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"text/template"

	"github.com/c4-project/c4t/internal/model/service/fuzzer"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/helper/errhelp"
)

const fuzzConfTemplate = `# AUTOGENERATED BY TESTER
fuzz {
{{ with .Machine }}## MACHINE SPECIFIC OVERRIDES ##
  # Set to number of cores in machine to prevent thrashing.
  set param cap.threads to {{ .Cores }}
{{ end -}}
{{- with .Config -}}## CONFIGURATION OVERRIDES ##
{{- range $k, $v := .Params }}
  {{ with $st := parseParam $k $v }}{{ $st }}{{ else }}# unsupported param "{{ $k }}": "{{ $v }}"{{ end }}
{{- end }}
{{ end -}}
}
`

func parseParam(key, val string) string {
	khd, ktl, ok := id.FromString(key).Uncons()
	if !ok {
		return ""
	}
	switch khd {
	case "int":
		return parseIntParam(ktl, val)
	case "bool":
		return parseBoolParam(ktl, val)
	case "action":
		return parseActionWeight(ktl, val)
	default:
		return ""
	}
}

func parseIntParam(key id.ID, val string) string {
	vint, err := strconv.Atoi(val)
	if err != nil {
		// ignore specific error
		return ""
	}
	return fmt.Sprintf("set param %s to %d", key, vint)
}

func parseBoolParam(key id.ID, val string) string {
	definites := map[string]bool{"on": true, "yes": true, "true": true, "off": false, "no": false, "false": false}
	b, ok := definites[val]
	if ok {
		return fmt.Sprintf("set flag %s to %s", key, boolToString(b))
	}
	return parseBoolRatioParam(key, val)
}

func parseBoolRatioParam(key id.ID, val string) string {
	var wins, losses int
	_, err := fmt.Sscanf(val, "%d:%d", &wins, &losses)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("set flag %s to ratio %d:%d", key, wins, losses)
}

func parseActionWeight(key id.ID, val string) string {
	vint, err := strconv.Atoi(val)
	if err != nil {
		// ignore specific error
		return ""
	}
	return fmt.Sprintf("action %s weight %d", key, vint)
}

func boolToString(b bool) string {
	// TODO(@MattWindsor91): is this redundant?
	if b {
		return "true"
	}
	return "false"
}

// WriteFuzzConf writes a fuzzer configuration based on j to w.
func WriteFuzzConf(w io.Writer, j fuzzer.Job) error {
	fmap := template.FuncMap{"parseParam": parseParam}
	t, err := template.New("fuzzConf").Funcs(fmap).Parse(fuzzConfTemplate)
	if err != nil {
		return err
	}
	return t.Execute(w, j)
}

// MakeFuzzConfFile creates a temporary file, then outputs WriteFuzzConf of j to it and returns the filepath.
// It is the caller's responsibility to delete the file.
func MakeFuzzConfFile(j fuzzer.Job) (string, error) {
	cf, err := ioutil.TempFile("", "act.*.conf")
	if err != nil {
		return "", fmt.Errorf("creating temporary fuzzer config file: %w", err)
	}

	werr := WriteFuzzConf(cf, j)
	cerr := cf.Close()

	return cf.Name(), errhelp.FirstError(werr, cerr)
}
