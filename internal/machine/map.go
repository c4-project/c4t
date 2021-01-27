// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"context"
	"errors"
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"
)

// ErrNoMachine occurs when we try to look up the compilers of a missing machine.
var ErrNoMachine = errors.New("no such machine")

// ConfigMap is a map from stringified IDs to machine configuration.
type ConfigMap map[string]Config

// Filter filters this Config's machines according to glob.
func (m ConfigMap) Filter(glob id.ID) (ConfigMap, error) {
	nm, err := id.MapGlob(m, glob)
	if err != nil {
		return nil, err
	}
	return nm.(ConfigMap), nil
}

// IDs gets a sorted slice of IDs present in this machine map.
// It returns an error if any of the configured machines have an invalid ID.
func (m ConfigMap) IDs() ([]id.ID, error) {
	return id.MapKeys(m)
}

// ListCompilers implements the compiler listing operation using a config.
func (m ConfigMap) ListCompilers(_ context.Context, mid id.ID) (map[string]compiler.Compiler, error) {
	mstr := mid.String()
	mach, ok := m[mstr]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoMachine, mstr)
	}
	return mach.Compilers, nil
}

// ObserveOn sends this map to obs as a series of machine observations.
func (m ConfigMap) ObserveOn(obs ...Observer) error {
	ids, err := m.IDs()
	if err != nil {
		return err
	}

	OnMachinesStart(len(ids), obs...)
	for i, n := range ids {
		OnMachinesRecord(i, Named{
			ID:      n,
			Machine: m[n.String()].Machine,
		}, obs...)
	}
	OnMachinesFinish(obs...)
	return nil
}
