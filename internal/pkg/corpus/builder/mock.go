// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package builder

import "github.com/MattWindsor91/act-tester/internal/pkg/model"

// MockObserver mocks Observer.
type MockObserver struct {
	// Manifest populates with the manifest when the observer observes OnStart.
	Manifest Manifest
	// Done sets to true when the observer observes OnFinish.
	Done bool
	// Adds tracks the add requests seen by the Observer.
	Adds map[string]struct{}
	// Compiles tracks the compile requests seen by the Observer.
	Compiles map[string][]model.ID
	// Harnesses tracks the harness requests seen by the Observer.
	Harnesses map[string][]model.ID
	// Runs tracks the run requests seen by the Observer.
	Runs map[string][]model.ID
}

// OnStart mocks the OnStart interface method.
func (t *MockObserver) OnStart(m Manifest) {
	t.Manifest = m
}

// OnRequest mocks the OnRequest interface method.
func (t *MockObserver) OnRequest(r Request) {
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

func (t *MockObserver) onCompile(sname string, cid model.ID) {
	addID(&t.Compiles, sname, cid)
}

func (t *MockObserver) onHarness(sname string, arch model.ID) {
	addID(&t.Harnesses, sname, arch)
}

func (t *MockObserver) onRun(sname string, cid model.ID) {
	addID(&t.Runs, sname, cid)
}

func addID(dest *map[string][]model.ID, key string, val model.ID) {
	if *dest == nil {
		*dest = map[string][]model.ID{}
	}
	(*dest)[key] = append((*dest)[key], val)
}

// OnFinish mocks the OnFinish interface method.
func (t *MockObserver) OnFinish() {
	t.Done = true
}
