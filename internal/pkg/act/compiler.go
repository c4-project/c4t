// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BinActCompiler is the name of the ACT compiler services binary.
const BinActCompiler = "act-compiler"

// ListCompilers queries ACT for a list of compilers satisfying f.
func (a *Runner) ListCompilers(ctx context.Context, f model.CompilerFilter) (map[string]map[string]model.Compiler, error) {
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer

	cmd := a.CommandContext(ctx, BinActCompiler, "list", sargs, f.ToArgv()...)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return ParseCompilerList(&obuf)
}

func (a *Runner) RunCompiler(ctx context.Context, c *model.NamedCompiler, infiles []string, outfile string, errw io.Writer) error {
	sargs := StandardArgs{Verbose: false}

	argv := runCompilerArgv(c.ID, infiles, outfile)
	cmd := a.CommandContext(ctx, BinActCompiler, "run", sargs, argv...)
	cmd.Stderr = errw

	return cmd.Run()
}

func runCompilerArgv(compiler model.ID, infiles []string, outfile string) []string {
	base := []string{"-compiler", compiler.String(), "-mode", "binary", "-o", outfile}
	return append(base, infiles...)
}

// ParseCompilerList parses a compiler list from the reader rd.
func ParseCompilerList(rd io.Reader) (map[string]map[string]model.Compiler, error) {
	compilers := make(map[string]map[string]model.Compiler)

	s := bufio.NewScanner(rd)
	for s.Scan() {
		mid, c, err := ParseCompiler(s.Bytes())
		if err != nil {
			return nil, err
		}

		if cerr := addCompiler(compilers, mid, c); cerr != nil {
			return nil, cerr
		}
	}

	return compilers, s.Err()
}

// addCompiler tries to add the compiler at machine CompilerID mid and compiler CompilerID mid to compilers.
// It fails if there is a duplicate compiler.
func addCompiler(compilers map[string]map[string]model.Compiler, mid model.ID, c model.NamedCompiler) error {
	ms := mid.String()
	if _, ok := compilers[ms]; !ok {
		compilers[ms] = make(map[string]model.Compiler)
	}

	cs := c.ID.String()
	if _, ok := compilers[ms][cs]; ok {
		return fmt.Errorf("duplicate compiler: machine=%s, compiler=%s", ms, cs)
	}

	compilers[ms][cs] = c.Compiler
	return nil
}

// ParseCompiler parses a single line from byte slice bs.
// It produces a machine CompilerID mid, a named compiler c, and/or an error.
func ParseCompiler(bs []byte) (mid model.ID, c model.NamedCompiler, err error) {
	s := bufio.NewScanner(bytes.NewReader(bs))
	s.Split(bufio.ScanWords)

	fields := []struct {
		name     string
		inserter func(string)
	}{
		{"machine CompilerID", func(s string) { mid = model.IDFromString(s) }},
		{"compiler CompilerID", func(s string) { c.ID = model.IDFromString(s) }},
		{"style", func(s string) { c.Style = model.IDFromString(s) }},
		{"arch", func(s string) { c.Arch = model.IDFromString(s) }},
		// enabled
	}

	for _, f := range fields {
		if !s.Scan() {
			return mid, c, CompilerFieldMissingError{
				line:  nil,
				field: f.name,
			}
		}
		f.inserter(s.Text())
	}

	return mid, c, nil
}

// CompilerFieldMissingError is an error caused when a compiler list line is missing an expected field.
type CompilerFieldMissingError struct {
	line  []byte
	field string
}

func (e CompilerFieldMissingError) Error() string {
	return fmt.Sprintf("no %s in compiler record %q", e.field, string(e.line))
}
