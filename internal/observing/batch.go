// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package observing

// Batch is a mixin struct for observer messages that represent 'batch runs'.
type Batch struct {
	// Kind is the kind of batch message being shown.
	Kind BatchKind `json:"kind"`

	// Num contains, in a start message, the number of items to expect in the batch.
	// In a step message, if set, it identifies the index of the item.
	Num int `json:"num"`
}

// BatchKind is the enumeration of kinds of batch message.
type BatchKind uint8

const (
	// A message that represents the start of a batch run.
	BatchStart BatchKind = iota
	// A message that represents a single step in a batch run.
	BatchStep
	// A message that represents the end of a batch run.
	BatchEnd
)

// NewBatchStart makes a batch start mixin with total count n.
func NewBatchStart(n int) Batch {
	return Batch{Kind: BatchStart, Num: n}
}

// NewBatchStep makes a batch step mixin with current index i.
func NewBatchStep(i int) Batch {
	return Batch{Kind: BatchStep, Num: i}
}

// NewBatchEnd makes a batch end mixin.
func NewBatchEnd() Batch {
	return Batch{Kind: BatchEnd}
}
