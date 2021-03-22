// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package tabulator

import (
	"io"
	"strings"
	"text/tabwriter"
)

// NewTab constructs a Tabulator that outputs elastic-tabbed human-readable output to w.
func NewTab(w io.Writer) Tabulator {
	return NewLineTabulator((*TabLineWriter)(tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)))
}

// TabLineWriter wraps a tabwriter.Writer to make it a LineWriter.
type TabLineWriter tabwriter.Writer

func (t *TabLineWriter) Write(cells []string) error {
	_, err := (*tabwriter.Writer)(t).Write([]byte(strings.Join(cells, "\t") + "\n"))
	return err
}

func (t *TabLineWriter) Flush() error {
	return (*tabwriter.Writer)(t).Flush()
}
