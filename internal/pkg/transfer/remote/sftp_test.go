// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/transfer/remote"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

type mockSFTP struct{ mock.Mock }

func (m *mockSFTP) Create(path string) (io.WriteCloser, error) {
	args := m.Called(path)
	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (m *mockSFTP) MkdirAll(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

type mockObserver struct{ mock.Mock }

func (m *mockObserver) OnCopyStart(nfiles int) {
	m.Called(nfiles)
}

func (m *mockObserver) OnCopy(src, dst string) {
	m.Called(src, dst)
}

func (m *mockObserver) OnCopyFinish() {
	m.Called()
}

type closeBuffer struct {
	bytes.Buffer
	// closed holds whether the buffer was closed.
	closed bool
}

func (c *closeBuffer) Close() error {
	c.closed = true
	return nil
}

// TestPutMapping tests PutMapping on a representative mapping.
func TestPutMapping(t *testing.T) {
	t.Parallel()

	// NB: the 'local' files here actually exist in the filesystem relative to this test.
	mapping := map[string]string{
		path.Join("remote", "bin", "a.out"):         path.Join("testdata", "sftp_test", "put1.txt"),
		path.Join("remote", "include", "foo.h"):     path.Join("testdata", "sftp_test", "put2.txt"),
		path.Join("remote", "src", "blah", "baz.c"): path.Join("testdata", "sftp_test", "put3.txt"),
	}

	var m mockSFTP

	for _, d := range []string{"bin", "include", path.Join("src", "blah")} {
		m.On("MkdirAll", path.Join("remote", d)).Return(nil).Once()
	}

	buffers := make(map[string]*closeBuffer, len(mapping))
	for r := range mapping {
		buffers[r] = new(closeBuffer)
		m.On("Create", r).Return(buffers[r], nil).Once()
	}

	var o mockObserver

	o.
		On("OnCopyStart", len(mapping)).Return().Once().
		On("OnCopyFinish").Return().Once()
	for r, l := range mapping {
		o.On("OnCopy", l, r).Return().Once()
	}

	err := remote.PutMapping(context.Background(), &m, &o, mapping)
	assert.NoError(t, err)

	if m.AssertExpectations(t) {
		for r, l := range mapping {
			bs, err := ioutil.ReadFile(l)
			assert.NoError(t, err, "reading local test file", l)
			assert.Equal(t, bs, buffers[r].Bytes(), "checking copy occurred from", l, "to", r)
			assert.True(t, buffers[r].closed, "buffer not closed for file", r)
		}
	}
	o.AssertExpectations(t)
}
