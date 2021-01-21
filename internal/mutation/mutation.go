// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package mutation contains support for mutation testing using c4t.
package mutation

// Mutant is an identifier for a particular mutant.
//
// Since we only support a mutation testing setups with integer mutant identifiers, this is just uint64 for now.
type Mutant = uint64

// EnvVar is the environment variable used for mutation testing.
//
// Some day, this might not be hard-coded.
const EnvVar = "C4_MUTATION"
