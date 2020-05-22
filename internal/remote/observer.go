// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

// CopyObserver is an interface for types that observe an SFTP file copy.
type CopyObserver interface {
	// OnCopyStart lets the observer know when a file copy (of nfiles files) is beginning.
	OnCopyStart(nfiles int)

	// OnCopy lets the observer know that a file copy (from path src to path dst) has happened.
	OnCopy(dst, src string)

	// OnCopyFinish lets the observer know when a file copy has finished.
	OnCopyFinish()
}

//go:generate mockery -name=CopyObserver

// OnCopyStart sends an OnCopyStart observation to multiple observers.
func OnCopyStart(nfiles int, cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopyStart(nfiles)
	}
}

// OnCopy sends an OnCopy observation to multiple observers.
func OnCopy(dst, src string, cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopy(dst, src)
	}
}

// OnCopyFinish sends an OnCopyFinish observation to multiple observers.
func OnCopyFinish(cos ...CopyObserver) {
	for _, o := range cos {
		o.OnCopyFinish()
	}
}
