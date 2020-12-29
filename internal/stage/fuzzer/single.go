// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"context"

	"github.com/c4-project/c4t/internal/model/service/fuzzer"
)

// SingleFuzzer represents types that can commune with a C litmus test fuzzer.
type SingleFuzzer interface {
	// Fuzz carries out the given fuzzing job.
	Fuzz(context.Context, fuzzer.Job) error
}

//go:generate mockery --name=SingleFuzzer

// NopFuzzer is a single-fuzzer that does nothing.
type NopFuzzer struct{}

// FuzzSingle does nothing, but pretends to fuzz a file.
func (n NopFuzzer) Fuzz(context.Context, fuzzer.Job) error {
	return nil
}
