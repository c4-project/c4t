// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import "context"

// SingleFuzzer represents types that can commune with a C litmus test fuzzer.
type SingleFuzzer interface {
	// FuzzSingle fuzzes the test at path inPath using the given seed,
	// outputting the new test to path outPath and the trace to tracePath.
	FuzzSingle(ctx context.Context, seed int32, inPath, outPath, tracePath string) error
}

// NopFuzzer is a single-fuzzer that does nothing.
type NopFuzzer struct{}

// FuzzSingle does nothing, but pretends to fuzz a file.
func (n NopFuzzer) FuzzSingle(_ context.Context, _ int32, _, _, _ string) error {
	return nil
}
