// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"io"

	"github.com/pkg/sftp"
)

// SFTPCopier wraps an sftp.Client to adapt it to the Copier interface.
type SFTPCopier sftp.Client

// Create wraps sftp.Client's Create in such a way as to implement Copier.
func (s *SFTPCopier) Create(path string) (io.WriteCloser, error) {
	return (*sftp.Client)(s).Create(path)
}

// Open wraps sftp.Client's Open in such a way as to implement Copier.
func (s *SFTPCopier) Open(path string) (io.ReadCloser, error) {
	return (*sftp.Client)(s).Open(path)
}

// MkdirAll wraps sftp.Client's MkdirAll in such a way as to implement Copier.
func (s *SFTPCopier) MkdirAll(dir string) error {
	return (*sftp.Client)(s).MkdirAll(dir)
}
