// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"io"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
)

// Option is the type of options for a pretty-printer.
type Option func(*Printer)

// Options combines the options os into a single option.
func Options(os ...Option) Option {
	return func(aw *Printer) {
		for _, o := range os {
			o(aw)
		}
	}
}

// WriteTo sets the printer's output to w.
func WriteTo(w io.Writer) Option {
	return func(aw *Printer) {
		aw.w = iohelp.EnsureWriter(w)
	}
}

// ShowOk sets whether the printer should show subjects in the 'ok' category, according to show.
func ShowOk(show bool) Option {
	return func(aw *Printer) {
		aw.ctx.ShowOk = show
	}
}

// ShowCompilers sets whether the printer should show compiler information, according to show.
func ShowCompilers(show bool) Option {
	return func(aw *Printer) {
		aw.ctx.ShowCompilers = show
	}
}

// ShowSubjects sets whether the printer should show subject breakdowns, according to show.
func ShowSubjects(show bool) Option {
	return func(aw *Printer) {
		aw.ctx.ShowSubjects = show
	}
}

// ShowPlanInfo sets whether the printer should show plan metadata, according to show.
func ShowPlanInfo(show bool) Option {
	return func(aw *Printer) {
		aw.ctx.ShowPlanInfo = show
	}
}

func ShowCompilerLogs(show bool) Option {
	return func(aw *Printer) {
		aw.ctx.ShowCompilerLogs = show
	}
}
