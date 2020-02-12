package model

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Compiler collects the test-relevant information about a compiler.
type Compiler struct {
	Service

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch Id
}

// ParseCompilerList parses a compiler list from the reader rd.
func ParseCompilerList(rd io.Reader) ([]*Compiler, error) {
	var cs []*Compiler

	s := bufio.NewScanner(rd)
	for s.Scan() {
		c, err := ParseCompiler(s.Bytes())
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}

	return cs, s.Err()
}

// ParseCompiler parses a single line from byte slice bs.
func ParseCompiler(bs []byte) (*Compiler, error) {
	s := bufio.NewScanner(bytes.NewReader(bs))
	s.Split(bufio.ScanWords)

	var c Compiler

	fields := []struct {
		name     string
		inserter func(*Compiler, string)
	}{
		{"machine Id", func(c *Compiler, s string) { m := IdFromString(s); c.MachineId = &m }},
		{"compiler Id", func(c *Compiler, s string) { c.Id = IdFromString(s) }},
		{"style", func(c *Compiler, s string) { c.Style = IdFromString(s) }},
		{"arch", func(c *Compiler, s string) { c.Arch = IdFromString(s) }},
		// enabled
	}

	for _, f := range fields {
		if !s.Scan() {
			return nil, CompilerFieldMissingError{
				line:  nil,
				field: f.name,
			}
		}
		f.inserter(&c, s.Text())
	}

	return &c, nil
}

// CompilerFieldMissingError is an error caused when a compiler list line is missing an expected field.
type CompilerFieldMissingError struct {
	line  []byte
	field string
}

func (e CompilerFieldMissingError) Error() string {
	return fmt.Sprintf("no %s in compiler record %q", e.field, string(e.line))
}
