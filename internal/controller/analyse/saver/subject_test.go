// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/observer"
	"github.com/MattWindsor91/act-tester/internal/model/filekind"
	"github.com/MattWindsor91/act-tester/internal/model/normaliser"
	"github.com/stretchr/testify/require"

	omocks "github.com/MattWindsor91/act-tester/internal/controller/analyse/observer/mocks"
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver"
	"github.com/MattWindsor91/act-tester/internal/controller/analyse/saver/mocks"
)

func TestArchiveSubject(t *testing.T) {
	var (
		ar  mocks.Archiver
		obs omocks.Observer
	)

	nm := normaliser.Map{
		normaliser.FileBin: normaliser.Entry{
			Original: "output.exe",
			Kind:     filekind.Bin,
			Loc:      filekind.InCompile,
		},
		normaliser.FileCompileLog: normaliser.Entry{
			Original: "output.log",
			Kind:     filekind.Log,
			Loc:      filekind.InCompile,
		},
		normaliser.FileOrigLitmus: normaliser.Entry{
			Original: "foo.litmus",
			Kind:     filekind.Litmus,
			Loc:      filekind.InOrig,
		},
		normaliser.FileFuzzLitmus: normaliser.Entry{
			Original: "foo_9.litmus",
			Kind:     filekind.Litmus,
			Loc:      filekind.InFuzz,
		},
		normaliser.FileFuzzTrace: normaliser.Entry{
			Original: "foo_9.trace",
			Kind:     filekind.Trace,
			Loc:      filekind.InFuzz,
		},
	}

	for n, m := range nm {
		ar.On("ArchiveFile", m.Original, n, m.Kind.ArchivePerm()).Return(nil).Once()
		orig := m.Original
		obs.On("OnArchive", mock.MatchedBy(func(s observer.ArchiveMessage) bool {
			return s.Kind == observer.ArchiveFileAdded && s.SubjectName == "foo_9" && s.File == orig
		})).Return().Once()
	}

	obs.On("OnArchive", mock.MatchedBy(func(s observer.ArchiveMessage) bool {
		return s.Kind == observer.ArchiveStart && s.Index == len(nm) && s.SubjectName == "foo_9" && s.File == "foo_9_saved"
	})).Return().Once()
	obs.On("OnArchive", mock.MatchedBy(func(s observer.ArchiveMessage) bool {
		return s.Kind == observer.ArchiveFinish && s.SubjectName == "foo_9"
	})).Return().Once()

	err := saver.ArchiveSubject(&ar, "foo_9", "foo_9_saved", nm, &obs)
	require.NoError(t, err, "ArchiveSubject")

	ar.AssertExpectations(t)
	obs.AssertExpectations(t)
}
