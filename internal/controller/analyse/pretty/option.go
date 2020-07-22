// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"io"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// Option is the type of options for a pretty-printer.
type Option func(*AnalysisWriter)

// Options combines the options os into a single option.
func Options(os ...Option) Option {
	return func(aw *AnalysisWriter) {
		for _, o := range os {
			o(aw)
		}
	}
}

// WriteTo sets the printer's output to w.
func WriteTo(w io.Writer) Option {
	return func(aw *AnalysisWriter) {
		aw.w = iohelp.EnsureWriter(w)
	}
}
