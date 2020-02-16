// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// BinaryRunner is the interface of things that can run compiled test binaries.
type BinaryRunner interface {
	// RunBinary runs the binary pointed to by bin, interpreting its results according to the backend spec backend.
	RunBinary(backend model.Backend, bin string)
}
