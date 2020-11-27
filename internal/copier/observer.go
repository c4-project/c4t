// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package copier

import "github.com/MattWindsor91/c4t/internal/observing"

// Observer is an interface for types that observe an SFTP file copy.
type Observer interface {
	// OnCopy sends a copy observation message.
	OnCopy(Message)
}

//go:generate mockery --name=Observer

// Message is the type of copy observation messages.
type Message struct {
	observing.Batch

	// Dst is the name of the destination file, if we're on a step.
	Dst string `json:"dst,omitempty"`

	// Src is the name of the source file, if we're on a step.
	Src string `json:"src,omitempty"`
}

// OnCopy sends an OnCopyStep observation to multiple observers.
func OnCopy(m Message, cos ...Observer) {
	for _, o := range cos {
		o.OnCopy(m)
	}
}

// OnCopyStart sends an OnCopyStart observation to multiple observers.
func OnCopyStart(nfiles int, cos ...Observer) {
	OnCopy(Message{Batch: observing.NewBatchStart(nfiles)}, cos...)
}

// OnCopyStep sends an OnCopyStep observation to multiple observers.
func OnCopyStep(i int, dst, src string, cos ...Observer) {
	OnCopy(Message{Batch: observing.NewBatchStep(i), Dst: dst, Src: src}, cos...)
}

// OnCopyEnd sends an OnCopyEnd observation to multiple observers.
func OnCopyEnd(cos ...Observer) {
	OnCopy(Message{Batch: observing.NewBatchEnd()}, cos...)
}
