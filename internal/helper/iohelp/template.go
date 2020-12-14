// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import "text/template"

// TemplateFromStrings parses a template from a series of named strings.
// It defines root first, under the name "root".
func TemplateFromStrings(root string, sub map[string]string) (*template.Template, error) {
	t, err := template.New("root").Parse(root)
	if err != nil {
		return nil, err
	}
	for n, ts := range sub {
		if t, err = t.New(n).Parse(ts); err != nil {
			return nil, err
		}
	}
	return t, nil
}
