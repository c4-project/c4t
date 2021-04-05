// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty

import (
	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/tabulator"

	"github.com/c4-project/c4t/internal/config"
)

// TabulateCompilers writes a full human-readable compiler table to w.
// It does not flush the table.
func TabulateCompilers(t tabulator.Tabulator, c *config.Config) error {
	setCompilerTableHeader(t)

	ms, err := c.Machines()
	if err != nil {
		return err
	}
	mids, err := ms.IDs()
	if err != nil {
		return err
	}
	for _, mid := range mids {
		if err := TabulateCompilersForMachine(t, mid, ms[mid]); err != nil {
			return err
		}
	}
	return t.Flush()
}

// TabulateCompilersForMachine writes a human-readable compiler table for machine mid/mc to w.
// It does not flush the table.
func TabulateCompilersForMachine(t tabulator.Tabulator, mid id.ID, mc machine.Config) error {
	// In case we're calling into this individually:
	setCompilerTableHeader(t)

	cs, err := mc.Compilers()
	if err != nil {
		return err
	}

	// Trying to force a deterministic order here.
	cids, err := id.MapKeys(cs)
	if err != nil {
		return err
	}

	// TODO(@MattWindsor91): deterministic order here?
	for _, cid := range cids {
		addCompilerRow(t, mid, cs[cid], cid)
	}

	return err
}

func addCompilerRow(t tabulator.Tabulator, mid id.ID, c compiler.Compiler, cid id.ID) {
	t.Cell(mid).Cell(cid).Cell(c.Style).Cell(c.Arch).Cell(status(&c)).EndRow()
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
