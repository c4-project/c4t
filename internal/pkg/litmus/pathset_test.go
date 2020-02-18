package litmus_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/litmus"
)

// ExamplePathset_MainCFile is a testable example for MainCFile.
func ExamplePathset_MainCFile() {
	fmt.Println((&litmus.Pathset{
		FileIn: "/home/foo/bar/baz.litmus",
		DirOut: "/tmp/scratch/",
	}).MainCFile())

	// Output:
	// /tmp/scratch/baz.c
}
