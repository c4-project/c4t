package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"

	"github.com/MattWindsor91/act-tester/internal/pkg/litmus"

	"github.com/MattWindsor91/act-tester/internal/pkg/ux"
)

const (
	usageCArch  = "C architecture to pass through to litmus"
	usageOutDir = "output directory for harness"
)

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		ux.LogTopError(err)
	}
}

// run contains the top-level logic of the command, passing through arguments args and stderr writer errw.
func run(args []string, errw io.Writer) error {
	cfg, err := parseArgs(args, errw)
	if err != nil {
		return err
	}

	return cfg.Run()
}

// parseArgs parses args into a litmus config, using errw as stderr.
func parseArgs(args []string, errw io.Writer) (*litmus.Litmus, error) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.SetOutput(errw)

	var act interop.ActRunner
	cfg := litmus.Litmus{
		Stat: act,
		Err:  errw,
	}

	fs.StringVar(&cfg.CArch, "carch", "", usageCArch)
	fs.StringVar(&cfg.OutDir, "o", "", usageOutDir)
	// TODO(@MattWindsor91): ActRunner flags

	if err := fs.Parse(args[1:]); err != nil {
		return nil, err
	}

	anons := fs.Args()
	if len(anons) != 1 {
		return nil, fmt.Errorf("expected precisely one anonymous argument; got %v", anons)
	}
	cfg.InFile = anons[0]
	return &cfg, nil
}
