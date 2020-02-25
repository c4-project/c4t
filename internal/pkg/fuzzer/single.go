// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// SingleFuzzer represents types that can commune with a C litmus test fuzzer.
type SingleFuzzer interface {
	// FuzzSingle fuzzes the test at path inPath using the given seed,
	// outputting files to the paths at outPaths.
	FuzzSingle(ctx context.Context, seed int32, inPath string, outPaths subject.FuzzFileset) error
}

// NopFuzzer is a single-fuzzer that does nothing.
type NopFuzzer struct{}

// FuzzSingle does nothing, but pretends to fuzz a file.
func (n NopFuzzer) FuzzSingle(_ context.Context, _ int32, _ string, _ subject.FuzzFileset) error {
	return nil
}
