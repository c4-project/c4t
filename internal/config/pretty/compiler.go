// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/tabulator"

	"github.com/c4-project/c4t/internal/config"
)

// TabulateCompilers writes a full human-readable compiler table to w.
// It does not flush the table.
func TabulateCompilers(t tabulator.Tabulator, c *config.Config) error {
	setCompilerTableHeader(t)

	ids, err := c.Machines.IDs()
	if err != nil {
		return err
	}
	for _, mid := range ids {
		TabulateCompilersForMachine(t, mid, c.Machines[mid.String()])
	}
	return t.Flush()
}

// TabulateCompilersForMachine writes a human-readable compiler table for machine mid/mc to w.
// It does not flush the table.
func TabulateCompilersForMachine(t tabulator.Tabulator, mid id.ID, mc machine.Config) {
	// In case we're calling into this individually:
	setCompilerTableHeader(t)

	// TODO(@MattWindsor91): deterministic order here?
	for cid, c := range mc.Compilers {
		addCompilerRow(t, mid, &c, cid)
	}
}

func addCompilerRow(t tabulator.Tabulator, mid id.ID, c *compiler.Compiler, cid string) {
	t.Cell(mid).Cell(cid).Cell(c.Style).Cell(c.Arch).Cell(status(c)).EndRow()
}

func status(c *compiler.Compiler) string {
	if c.Disabled {
		return "off"
	}
	return "on"
}

func setCompilerTableHeader(t tabulator.Tabulator) {
	t.Header("Machine", "ID", "Style", "Arch", "Status")
}
