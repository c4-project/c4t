// Package litmus implements a wrapper over the herdtools7 'litmus7' program.
// This wrapper deals with various corner cases.

package litmus

import (
	"errors"
	"io"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
)

var (
	// ErrConfigNil occurs when Run tries to run on a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrStatNil occurs when the config has a nil statistics dumper.
	ErrStatNil = errors.New("config stat dumper nil")

	// ErrNoCArch occurs when the output directory is empty.
	ErrNoCArch = errors.New("need carch")

	// ErrNoInFile occurs when the input file is empty.
	ErrNoInFile = errors.New("need input file")

	// ErrNoOutDir occurs when the output directory is empty.
	ErrNoOutDir = errors.New("need output directory")
)

// Litmus is the configuration required to run the litmus shim.
type Litmus struct {
	// Stat extracts statistics from litmus files.
	// These statistics then switch on various fixes.
	Stat interop.StatDumper

	// CArch is the architecture that the litmus shim should target.
	// It corresponds to Litmus's 'carch' argument.
	CArch string

	// InFile is the path to the input test file.
	InFile string

	// OutDir is the path to the output directory.
	OutDir string

	// Fixset is the set of enabled fixes.
	// It is part of the config to allow the forcing of fixes that the shim would otherwise deem unnecessary.
	Fixset Fixset

	// Err is the writer to which stderr output should be written.
	Err io.Writer
}

// Run runs the litmus wrapper according to the configuration c.
func (c *Litmus) Run() error {
	if err := c.check(); err != nil {
		return err
	}

	if err := c.probeFixes(); err != nil {
		return err
	}

	if err := c.runLitmus(); err != nil {
		return err
	}

	return nil
}

// check checks that the configuration makes sense.
func (c *Litmus) check() error {
	if c == nil {
		return ErrConfigNil
	}
	if c.Stat == nil {
		return ErrStatNil
	}
	if c.CArch == "" {
		return ErrNoCArch
	}
	if c.InFile == "" {
		return ErrNoInFile
	}
	if c.OutDir == "" {
		return ErrNoOutDir
	}
	return nil
}

// probeFixes checks to see if there are any fixes needed for the input.
func (c *Litmus) probeFixes() error {
	var s interop.Statset
	if err := c.Stat.DumpStats(&s, c.InFile); err != nil {
		return err
	}
	c.Fixset.PopulateFromStats(&s)
	return c.Fixset.Dump(c.Err)
}

// runLitmus actually runs Litmus.
func (c *Litmus) runLitmus() error {
	cmd := exec.Command("litmus7", c.litmusArgs()...)

	cmd.Stderr = c.Err
	return cmd.Run()
}

// litmusArgs works out the argument vector for Litmus.
func (c *Litmus) litmusArgs() []string {
	args := c.Fixset.Args()
	args = append(args, "-carch", c.CArch, "-c11", "true", "-o", c.OutDir, c.InFile)
	return args
}
