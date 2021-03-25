// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/ux/stdflag"

	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/config"

	app "github.com/c4-project/c4t/internal/app/config"
	"github.com/stretchr/testify/require"
)

// TestApp_printGlobalPath tests the config app's ability to dump the global path of the config file.
func TestApp_printGlobalPath(t *testing.T) {
	var buf bytes.Buffer

	// TODO(@MattWindsor91): we should really compare this err to the one from GlobalPath.
	err := app.App(&buf, io.Discard).Run([]string{app.Name, "-" + app.FlagPrintGlobalPath})
	require.NoError(t, err, "error while running app")

	gp, err := config.GlobalFile()
	require.NoError(t, err, "error while getting global config file")

	// We can't actually assert that the global file path _is_ something in particular, as it'll depend on OS
	require.Equal(t, gp, strings.TrimSpace(buf.String()), "should have printed global file to stdout")
}

// TestApp_compilerTables tests the output of 'c4t-config -print-compilers' given various input files.
func TestApp_compilerTables(t *testing.T) {
	t.Parallel()

	testhelp.TestFilesOfExt(t, path.Join("testdata", "compilers"), ".toml", func(t *testing.T, name, path string) {
		t.Helper()
		t.Parallel()
		a := assert.New(t)

		var buf bytes.Buffer

		if err := app.App(&buf, io.Discard).Run([]string{
			app.Name,
			"-" + app.FlagPrintCompilers,
			"-" + stdflag.FlagConfigFile,
			path,
		}); !a.NoError(err) {
			return
		}

		got := buf.String()
		wbytes, err := os.ReadFile(filepath.Join("testdata", "compilers", name+".txt"))
		if !a.NoError(err) {
			return
		}

		a.Equalf(string(wbytes), got, "compiler output not equal for %q", name)
	})
}
