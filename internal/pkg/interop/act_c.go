package interop

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BinActC is the name of the ACT C services binary.
const BinActC = "act-c"

// ErrSubjectNil is an error returned if one calls ProbeSubject(nil).
var ErrSubjectNil = errors.New("subject pointer is nil")

// ProbeSubject populates subject with information gleaned from investigating its litmus file.
func (a ActRunner) ProbeSubject(subject *model.Subject) error {
	if subject == nil {
		return ErrSubjectNil
	}

	var obuf bytes.Buffer
	sargs := StandardArgs{Verbose: false}

	if err := a.Run(BinActC, nil, &obuf, nil, sargs, "dump-header", subject.Litmus); err != nil {
		return fmt.Errorf("ACT dump-header on %s failed: %w", subject.Litmus, err)
	}
	hdr, err := ReadHeader(&obuf)
	if err != nil {
		return fmt.Errorf("header read on %s failed: %w", subject.Litmus, err)
	}

	probeSubjectFromHeader(subject, hdr)
	return nil
}

func probeSubjectFromHeader(subject *model.Subject, h *Header) {
	subject.Name = h.Name
	// TODO(@MattWindsor91): number of threads
}
