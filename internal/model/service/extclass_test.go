// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestExtClass_ProbeByVersionCommand tests ProbeByVersionCommand by using a dummy service runner.
func TestExtClass_ProbeByVersionCommand(t *testing.T) {
	t.Parallel()

	var sr mocks.Runner
	sr.Test(t)

	// The command given in the default run arguments will look like it doesn't exist, so shouldn't appear in the final
	// dict.  The alternative commands are the ones that will look like they exist, and should give a final dictionary
	// equal to `want`.

	ec := service.ExtClass{
		DefaultRunInfo: service.RunInfo{Cmd: "nothere", Args: []string{"-blah"}},
		AltCommands:    []string{"foo", "bar", "baz"},
	}
	argsValid := func(args []string) bool {
		// Should preserve -blah from default run info and add -version from the call.
		return len(args) == 2 && args[0] == "-blah" && args[1] == "-version"
	}

	want := map[string]string{"foo": "v1.01", "bar": "356", "baz": "4.10.1998"}

	// Need to mock the stdout redirection.  We assume we don't parallelise the search here, this may not be wise.
	var buf io.Writer
	sr.On("WithStdout", mock.Anything).Run(func(args mock.Arguments) {
		buf = args.Get(0).(io.Writer)
	}).Return(&sr).Times(len(want) + 1)

	// Set up the command that doesn't exist.
	sr.On("Run", mock.Anything, mock.MatchedBy(func(info service.RunInfo) bool {
		return info.Cmd == ec.DefaultRunInfo.Cmd
	})).Return(errors.New("not here")).Once()

	// Set up the commands that do exist.
	for k, v := range want {
		k, v := k, v

		sr.On("Run", mock.Anything, mock.MatchedBy(func(info service.RunInfo) bool {
			return info.Cmd == k && argsValid(info.Args)
		})).Run(func(_ mock.Arguments) {
			_, _ = buf.Write([]byte(v))
		}).Return(nil).Once()
	}

	// Now we can actually do the test.
	assert.Equal(t, want, ec.ProbeByVersionCommand(context.Background(), &sr, "-version"))

	sr.AssertExpectations(t)
}
