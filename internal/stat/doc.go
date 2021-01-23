// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package stat implements c4t's persistent statistics support.
//
// This includes the models for statistics collection (Set, MachineSet, etc), and the Persister, a director observer
// that tracks statistics by persisting them to a JSON file.
package stat
