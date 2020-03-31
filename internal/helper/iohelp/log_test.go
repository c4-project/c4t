// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
)

// TestEnsureLog tests various properties of EnsureLog.
func TestEnsureLog(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		l := iohelp.EnsureLog(nil)
		if l == nil {
			t.Fatal("EnsureLog(nil) came back nil")
		}
		if l.Writer() != ioutil.Discard {
			t.Fatal("EnsureLog(nil) uses strange writer:", l.Writer())
		}
	})
	t.Run("non-nil", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		ol := log.New(&buf, "[test] ", log.LstdFlags)
		l := iohelp.EnsureLog(ol)
		if l != ol {
			t.Fatal("EnsureLog(l) didn't return pointer unchanged")
		}
	})
}
