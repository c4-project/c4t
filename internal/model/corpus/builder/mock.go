// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// MockObserver mocks Observer.
type MockObserver struct {
	// Manifest populates with the manifest when the observer observes OnBuildStart.
	Manifest Manifest
	// Done sets to true when the observer observes OnBuildFinish.
	Done bool
	// Adds tracks the add requests seen by the Observer.
	Adds map[string]struct{}
	// Compiles tracks the compile requests seen by the Observer.
	Compiles map[string][]id.ID
	// Harnesses tracks the harness requests seen by the Observer.
	Harnesses map[string][]id.ID
	// Runs tracks the run requests seen by the Observer.
	Runs map[string][]id.ID
}

// OnBuildStart mocks the OnBuildStart interface method.
func (t *MockObserver) OnBuildStart(m Manifest) {
	t.Manifest = m
}

// OnBuildRequest mocks the OnBuildRequest interface method.
func (t *MockObserver) OnBuildRequest(r Request) {
	switch {
	case r.Add != nil:
		t.onAdd(r.Name)
	case r.Compile != nil:
		t.onCompile(r.Name, r.Compile.CompilerID)
	case r.Harness != nil:
		t.onHarness(r.Name, r.Harness.Arch)
	case r.Run != nil:
		t.onRun(r.Name, r.Run.CompilerID)
	}
}

func (t *MockObserver) onAdd(sname string) {
	if t.Adds == nil {
		t.Adds = map[string]struct{}{}
	}
	t.Adds[sname] = struct{}{}
}

func (t *MockObserver) onCompile(sname string, cid id.ID) {
	addID(&t.Compiles, sname, cid)
}

func (t *MockObserver) onHarness(sname string, arch id.ID) {
	addID(&t.Harnesses, sname, arch)
}

func (t *MockObserver) onRun(sname string, cid id.ID) {
	addID(&t.Runs, sname, cid)
}

func addID(dest *map[string][]id.ID, key string, val id.ID) {
	if *dest == nil {
		*dest = map[string][]id.ID{}
	}
	(*dest)[key] = append((*dest)[key], val)
}

// OnBuildFinish mocks the OnBuildFinish interface method.
func (t *MockObserver) OnBuildFinish() {
	t.Done = true
}
