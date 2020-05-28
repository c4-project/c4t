// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package runner

// InvocationGetter is the interface of types that tell the invoker how to invoke the machine node.
type InvocationGetter interface {
	// MachBin retrieves the default binary name for the node.
	MachBin() string
	// MachArgs computes the argument set for invoking the node binary.
	MachArgs(dir string) []string
}

// Invocation gets the invocation for the local-machine binary as a string list.
func Invocation(i InvocationGetter, dir string) []string {
	return append([]string{i.MachBin()}, i.MachArgs(dir)...)
}
