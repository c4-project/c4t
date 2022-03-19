// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"context"
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"
)

// ErrNoMachine occurs when we try to look up the compilers of a missing machine.
var ErrNoMachine = errors.New("no such machine")

// ConfigMap is a map from IDs to machine configuration.
type ConfigMap map[id.ID]Config

// Filter filters this Config's machines according to glob.
func (m ConfigMap) Filter(glob id.ID) (ConfigMap, error) {
	nm, err := id.MapGlob(m, glob)
	if err != nil {
		return nil, err
	}
	return nm, nil
}

// IDs gets a sorted slice of IDs present in this machine map.
// It returns an error if any of the configured machines have an invalid ID.
func (m ConfigMap) IDs() []id.ID {
	return id.MapKeys(m)
}

// ListCompilers implements the compiler listing operation using a config.
func (m ConfigMap) ListCompilers(_ context.Context, machine id.ID) (map[id.ID]compiler.Compiler, error) {
	mach, ok := m[machine]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoMachine, machine)
	}
	return mach.Compilers()
}

// ObserveOn sends this map to obs as a series of machine observations.
func (m ConfigMap) ObserveOn(obs ...Observer) error {
	ids := m.IDs()

	OnMachinesStart(len(ids), obs...)
	for i, n := range ids {
		OnMachinesRecord(i, Named{
			ID:      n,
			Machine: m[n].Machine,
		}, obs...)
	}
	OnMachinesFinish(obs...)
	return nil
}
