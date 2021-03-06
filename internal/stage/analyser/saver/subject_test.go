// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package saver_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/subject/normpath"

	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/model/filekind"
	"github.com/c4-project/c4t/internal/subject/normaliser"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"
	"github.com/c4-project/c4t/internal/stage/analyser/saver/mocks"
)

// TestArchiveSubject tests ArchiveSubject.
func TestArchiveSubject(t *testing.T) {
	var (
		ar  mocks.Archiver
		obs mocks.Observer
	)
	ar.Test(t)
	obs.Test(t)

	nm := normaliser.Map{
		normpath.FileBin: normaliser.Entry{
			Original: "output.exe",
			Kind:     filekind.Bin,
			Loc:      filekind.InCompile,
		},
		normpath.FileCompileLog: normaliser.Entry{
			Original: "output.log",
			Kind:     filekind.Log,
			Loc:      filekind.InCompile,
		},
		normpath.FileOrigLitmus: normaliser.Entry{
			Original: "foo.litmus",
			Kind:     filekind.Litmus,
			Loc:      filekind.InOrig,
		},
		normpath.FileFuzzLitmus: normaliser.Entry{
			Original: "foo_9.litmus",
			Kind:     filekind.Litmus,
			Loc:      filekind.InFuzz,
		},
		normpath.FileFuzzTrace: normaliser.Entry{
			Original: "foo_9.trace",
			Kind:     filekind.Trace,
			Loc:      filekind.InFuzz,
		},
	}

	for n, m := range nm {
		ar.On("ArchiveFile", m.Original, n, m.Kind.ArchivePerm()).Return(nil).Once()
		orig := m.Original
		obs.On("OnArchive", mock.MatchedBy(func(s saver.ArchiveMessage) bool {
			return s.Kind == saver.ArchiveFileAdded && s.SubjectName == "foo_9" && s.File == orig
		})).Return().Once()
	}

	obs.On("OnArchive", mock.MatchedBy(func(s saver.ArchiveMessage) bool {
		return s.Kind == saver.ArchiveStart && s.Index == len(nm) && s.SubjectName == "foo_9" && s.File == "foo_9_saved"
	})).Return().Once()
	obs.On("OnArchive", mock.MatchedBy(func(s saver.ArchiveMessage) bool {
		return s.Kind == saver.ArchiveFinish && s.SubjectName == "foo_9"
	})).Return().Once()

	err := saver.ArchiveSubject(&ar, "foo_9", "foo_9_saved", nm, &obs)
	require.NoError(t, err, "ArchiveSubject")

	ar.AssertExpectations(t)
	obs.AssertExpectations(t)
}
