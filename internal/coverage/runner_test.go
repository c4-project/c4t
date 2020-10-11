// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"fmt"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/coverage"
	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/subject"
)

// ExampleRunnerContext_ExpandArgs is a runnable example for RunnerContext.ExpandArgs.
func ExampleRunnerContext_ExpandArgs() {
	rc := coverage.RunnerContext{
		Seed:        8675309,
		BucketDir:   "bucket1,1",
		NumInBucket: 42,
		Input:       subject.NewOrPanic(litmus.New("foo/bar.litmus")),
	}
	args := rc.ExpandArgs("-seed", "${seed}", "-o", "${outputDir}/${i}.c", "${input}")
	for _, arg := range args {
		fmt.Println(filepath.ToSlash(arg))
	}

	// Output:
	// -seed
	// 8675309
	// -o
	// bucket1,1/42.c
	// foo/bar.litmus
}
