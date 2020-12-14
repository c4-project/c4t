// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import "text/template"

// ParseTemplateStrings adds to t the template bindings in sub.
func ParseTemplateStrings(t *template.Template, sub map[string]string) (*template.Template, error) {
	var err error
	for n, ts := range sub {
		if t, err = t.New(n).Parse(ts); err != nil {
			return nil, err
		}
	}
	return t, nil
}
