// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"io"

	"github.com/c4-project/c4t/internal/stat"

	"github.com/c4-project/c4t/internal/helper/iohelp"
)

// Option is the type of options for a pretty-printer.
type Option func(*Printer)

// Options combines the options os into a single option.
func Options(os ...Option) Option {
	return func(p *Printer) {
		for _, o := range os {
			o(p)
		}
	}
}

// WriteTo sets the printer's output to w.
func WriteTo(w io.Writer) Option {
	return func(p *Printer) {
		p.w = iohelp.EnsureWriter(w)
	}
}

// ShowMutants enables mutant showing under filter f.
// If f is nil, mutant showing is disabled.
func ShowMutants(f stat.MutantFilter) Option {
	return func(p *Printer) {
		p.ctx.MutantFilter = f
	}
}

// UseTotals determines whether we show totals rather than session statistics.
func UseTotals(on bool) Option {
	return func(p *Printer) {
		p.ctx.UseTotals = on
	}
}
