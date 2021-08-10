// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package database contains database functionality for c4t.
package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// DefaultPath is the default relative slashpath to the database.
const DefaultPath = "c4t.db"

// Open opens a SQLite database connection to the file on the given path.
func Open(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}
