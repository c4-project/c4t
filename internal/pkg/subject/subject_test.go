package subject

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ExampleSubject_CompileResult is a testable example for CompileResult.
func ExampleSubject_CompileResult() {
	s := Subject{Compiles: map[string]CompileResult{
		"localhost:gcc":   {Success: true, Files: CompileFileset{Bin: "a.out", Log: "gcc.log"}},
		"spikemuth:clang": {Success: false, Files: CompileFileset{Bin: "a.out", Log: "clang.log"}},
	}}
	lps, _ := s.CompileResult(model.MachQualID{MachineID: model.IDFromString("localhost"), ID: model.IDFromString("gcc")})
	sps, _ := s.CompileResult(model.MachQualID{MachineID: model.IDFromString("spikemuth"), ID: model.IDFromString("clang")})

	fmt.Println("localhost:", lps.Success, lps.Files.Bin, lps.Files.Log)
	fmt.Println("spikemuth:", sps.Success, sps.Files.Bin, sps.Files.Log)

	// Output:
	// localhost: true a.out gcc.log
	// spikemuth: false a.out clang.log
}

// ExampleSubject_Harness is a testable example for Harness.
func ExampleSubject_Harness() {
	s := Subject{Harnesses: map[string]Harness{
		"localhost:x86.64": {Dir: "foo", Files: []string{"bar", "baz"}},
		"spikemuth:arm":    {Dir: "foobar", Files: []string{"barbaz"}},
	}}
	lps, _ := s.Harness(model.MachQualID{MachineID: model.IDFromString("localhost"), ID: model.ArchX8664})
	sps, _ := s.Harness(model.MachQualID{MachineID: model.IDFromString("spikemuth"), ID: model.ArchArm})

	for _, l := range lps.Files {
		fmt.Println(l)
	}
	for _, s := range sps.Files {
		fmt.Println(s)
	}

	// Output:
	// bar
	// baz
	// barbaz
}

// TestSubject_CompileResult_Missing checks that trying to get a harness path for a missing machine/emits pair triggers
// the appropriate error.
func TestSubject_CompileResult_Missing(t *testing.T) {
	var s Subject
	_, err := s.CompileResult(model.MachQualID{
		MachineID: model.IDFromString("localhost"),
		ID:        model.IDFromString("gcc"),
	})
	testhelp.ExpectErrorIs(t, err, ErrMissingCompile, "missing compile result path")
}

// TestSubject_AddCompileResult checks that AddCompileResult is working properly.
func TestSubject_AddCompileResult(t *testing.T) {
	var s Subject
	c := CompileResult{
		Success: true,
		Files: CompileFileset{
			Bin: "a.out",
			Log: "gcc.log",
		},
	}

	mcomp := model.MachQualID{
		MachineID: model.IDFromString("localhost"),
		ID:        model.IDFromString("gcc"),
	}

	t.Run("initial-add", func(t *testing.T) {
		if err := s.AddCompileResult(mcomp, c); err != nil {
			t.Fatalf("err when adding compile to empty subject: %v", err)
		}
	})
	t.Run("add-get", func(t *testing.T) {
		c2, err := s.CompileResult(mcomp)
		if err != nil {
			t.Fatalf("err when getting added compile: %v", err)
		}
		if !reflect.DeepEqual(c2, c) {
			t.Fatalf("added compile (%v) came back wrong (%v)", c2, c)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		err := s.AddCompileResult(mcomp, CompileResult{})
		testhelp.ExpectErrorIs(t, err, ErrDuplicateCompile, "adding compile twice")
	})
}

// TestSubject_Harness_Missing checks that trying to get a harness path for a missing machine/emits pair triggers
// the appropriate error.
func TestSubject_Harness_Missing(t *testing.T) {
	var s Subject
	_, err := s.Harness(model.MachQualID{MachineID: model.IDFromString("localhost"), ID: model.IDFromString("x86.64")})
	testhelp.ExpectErrorIs(t, err, ErrMissingHarness, "missing harness path")
}

// TestSubject_AddHarness checks that AddHarness is working properly.
func TestSubject_AddHarness(t *testing.T) {
	var s Subject
	h := Harness{
		Dir:   "foo",
		Files: []string{"bar", "baz"},
	}

	march := model.MachQualID{
		MachineID: model.IDFromString("localhost"),
		ID:        model.ArchX8664,
	}

	t.Run("initial-add", func(t *testing.T) {
		if err := s.AddHarness(march, h); err != nil {
			t.Fatalf("err when adding harness to empty subject: %v", err)
		}
	})
	t.Run("add-get", func(t *testing.T) {
		h2, err := s.Harness(march)
		if err != nil {
			t.Fatalf("err when getting added harness: %v", err)
		}
		if !reflect.DeepEqual(h2, h) {
			t.Fatalf("added harness (%v) came back wrong (%v)", h2, h)
		}
	})
	t.Run("add-dupe", func(t *testing.T) {
		err := s.AddHarness(march, Harness{})
		testhelp.ExpectErrorIs(t, err, ErrDuplicateHarness, "adding harness twice")
	})
}
