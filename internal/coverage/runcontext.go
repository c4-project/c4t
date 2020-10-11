// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// RunContext is the type of state provided to a coverage runner.
type RunContext struct {
	// Seed is the seed to use to drive any random parts of the coverage runner.
	Seed int32
	// BucketDir is the filepath to the bucket directory into which the coverage runner should output its recipe.
	BucketDir string
	// NumInBucket is the index of this single instance in its bucket.
	NumInBucket int
	// Input points to an input subject for the coverage runner, if any are available.
	Input *subject.Subject
}

// inputPath tries to get the filepath to the currently available input's litmus test.
func (r RunContext) inputPath() (string, error) {
	if r.Input == nil {
		return "", ErrNoInput
	}
	l, err := r.Input.BestLitmus()
	if err != nil {
		return "", err
	}
	return filepath.Clean(l.Path), nil
}

func (r RunContext) inputPathOrEmpty() string {
	in, err := r.inputPath()
	if err != nil {
		return ""
	}
	return in
}

// ExpandArgs expands various special identifiers in args to parts of the runner context.
func (r RunContext) ExpandArgs(arg ...string) []string {
	replacer := strings.NewReplacer(
		"${seed}", strconv.Itoa(int(r.Seed)),
		"${input}", r.inputPathOrEmpty(),
		"${outputDir}", r.BucketDir,
		"${i}", strconv.Itoa(r.NumInBucket),
	)
	nargs := make([]string, len(arg))
	for i, a := range arg {
		nargs[i] = replacer.Replace(a)
	}
	return nargs
}

// LiftOutDir is a suggested output directory for backend lifts in this runner context.
func (r RunContext) LiftOutDir() string {
	return r.out("_harness")
}

// OutLitmus is a suggested filename for Litmus outputs of this runner context.
func (r RunContext) OutLitmus() string {
	return r.out(".litmus")
}

func (r RunContext) out(suffix string) string {
	return filepath.Join(r.BucketDir, fmt.Sprintf("%d%s", r.NumInBucket, suffix))
}
