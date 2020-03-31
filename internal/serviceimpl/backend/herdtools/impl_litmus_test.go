// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// ExampleLitmus_Args is a testable example for Args.
func ExampleLitmus_Args() {
	j := job.Harness{
		Arch:   id.ArchX8664,
		InFile: "in.litmus",
		OutDir: "out",
	}
	r := service.RunInfo{
		Cmd:  "litmus7",
		Args: []string{"-v"},
	}

	args64, _ := Litmus{}.Args(j, r)
	fmt.Print("64-bit:")
	for _, arg := range args64 {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()

	// 32-bit x86 maps to a different Litmus architecture:
	j.Arch = id.ArchX86
	args32, _ := Litmus{}.Args(j, r)
	fmt.Print("32-bit:")
	for _, arg := range args32 {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()

	// Output:
	// 64-bit: -o out -carch X86_64 -c11 true -v in.litmus
	// 32-bit: -o out -carch X86 -c11 true -v in.litmus
}
