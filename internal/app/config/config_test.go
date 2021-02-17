// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

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

	require.Equal(t, gp, strings.TrimSpace(buf.String()), "should have printed global file to stdout")
}
