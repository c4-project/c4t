// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path"
)

//go:embed sql
var schemata embed.FS

// Create creates the c4t schemata in the database db.
//
// While the driver can be any sql.DB compatible driver, the sql are written
// with sqlite in mind.
//
// Create takes an observer function to which the name of any DDL file being read is
// passed (can be nil).
func Create(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := createInTransaction(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func createInTransaction(ctx context.Context, tx *sql.Tx) error {
	bs, err := schemata.ReadFile(path.Join("sql", "schema.sql"))
	if err != nil {
		return fmt.Errorf("couldn't read schema: %w", err)
	}
	_, err = tx.ExecContext(ctx, string(bs))
	if err != nil {
		return fmt.Errorf("couldn't apply schema: %w", err)
	}
	return nil
}
