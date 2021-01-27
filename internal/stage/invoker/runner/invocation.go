// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package runner

// InvocationGetter is the interface of types that tell the invoker how to invoke the machine node.
type InvocationGetter interface {
	// MachBin retrieves the default binary name for the node.
	MachBin() string
	// MachArgs computes the argument set for invoking the node binary.
	MachArgs(dir string) []string
}
