// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"fmt"

	"github.com/MattWindsor91/c4t/internal/model/id"
	"github.com/MattWindsor91/c4t/internal/model/litmus"
)

// ExampleLitmus_IsC is a testable example for Litmus.IsC.
func ExampleLitmus_IsC() {
	foo := litmus.NewOrPanic("foo.litmus", litmus.WithArch(id.ArchC))
	fmt.Println("C:  ", foo.IsC())

	bar := litmus.NewOrPanic("bar.litmus", litmus.WithArch(id.ArchC.Join(id.FromString("11"))))
	fmt.Println("C11:", bar.IsC())

	baz := litmus.NewOrPanic("baz.litmus", litmus.WithArch(id.ArchArm))
	fmt.Println("Arm:", baz.IsC())

	// Output:
	// C:   true
	// C11: true
	// Arm: false
}
