package model

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Compiler collects the test-relevant information about a compiler.
type Compiler struct {
	// Style is the declared style of the backend.
	Style ID `toml:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch ID `toml:"arch"`
}

// ParseCompilerList parses a compiler list from the reader rd.
func ParseCompilerList(rd io.Reader) (map[string]map[string]Compiler, error) {
	compilers := make(map[string]map[string]Compiler)

	s := bufio.NewScanner(rd)
	for s.Scan() {
		mid, cid, c, err := ParseCompiler(s.Bytes())
		if err != nil {
			return nil, err
		}

		if cerr := addCompiler(compilers, mid, cid, c); cerr != nil {
			return nil, cerr
		}
	}

	return compilers, s.Err()
}

// addCompiler tries to add the compiler at machine ID mid and compiler ID mid to compilers.
// It fails if there is a duplicate compiler.
func addCompiler(compilers map[string]map[string]Compiler, mid, cid ID, c Compiler) error {
	ms := mid.String()
	if _, ok := compilers[ms]; !ok {
		compilers[ms] = make(map[string]Compiler)
	}

	cs := cid.String()
	if _, ok := compilers[ms][cs]; ok {
		return fmt.Errorf("duplicate compiler: machine=%s, compiler=%s", ms, cs)
	}

	compilers[ms][cs] = c
	return nil
}

// ParseCompiler parses a single line from byte slice bs.
// It produces a machine ID mid, a compiler ID cid, a compiler pointer, and/or an error.
func ParseCompiler(bs []byte) (mid, cid ID, c Compiler, err error) {
	s := bufio.NewScanner(bytes.NewReader(bs))
	s.Split(bufio.ScanWords)

	fields := []struct {
		name     string
		inserter func(string)
	}{
		{"machine ID", func(s string) { mid = IDFromString(s) }},
		{"compiler ID", func(s string) { cid = IDFromString(s) }},
		{"style", func(s string) { c.Style = IDFromString(s) }},
		{"arch", func(s string) { c.Arch = IDFromString(s) }},
		// enabled
	}

	for _, f := range fields {
		if !s.Scan() {
			return mid, cid, c, CompilerFieldMissingError{
				line:  nil,
				field: f.name,
			}
		}
		f.inserter(s.Text())
	}

	return mid, cid, c, nil
}

// CompilerFieldMissingError is an error caused when a compiler list line is missing an expected field.
type CompilerFieldMissingError struct {
	line  []byte
	field string
}

func (e CompilerFieldMissingError) Error() string {
	return fmt.Sprintf("no %s in compiler record %q", e.field, string(e.line))
}
