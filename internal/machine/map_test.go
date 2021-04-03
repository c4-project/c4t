// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/machine/mocks"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/machine"
)

// ExampleConfigMap_IDs is a runnable example for IDs.
func ExampleConfigMap_IDs() {
	cm := machine.ConfigMap{
		"localhost": machine.Config{Machine: machine.Machine{Cores: 3}},
		"bar":       machine.Config{Machine: machine.Machine{Cores: 1}},
		"foo.bar":   machine.Config{Machine: machine.Machine{Cores: 2}},
	}
	ids, _ := cm.IDs()
	for _, n := range ids {
		fmt.Println(n)
	}

	// Output:
	// bar
	// foo.bar
	// localhost
}

// ExampleConfigMap_Filter is a runnable example for Filter.
func ExampleConfigMap_Filter() {
	cm := machine.ConfigMap{
		"bar":         machine.Config{Machine: machine.Machine{Cores: 1}},
		"foo.bar":     machine.Config{Machine: machine.Machine{Cores: 2}},
		"foo.bar.baz": machine.Config{Machine: machine.Machine{Cores: 3}},
		"foo.baz":     machine.Config{Machine: machine.Machine{Cores: 4}},
	}
	cm, _ = cm.Filter(id.FromString("foo.*.baz"))
	ids, _ := cm.IDs()
	for _, n := range ids {
		fmt.Println(n)
	}

	// Output:
	// foo.bar.baz
	// foo.baz
}

func TestConfigMap_ObserveOn(t *testing.T) {
	t.Parallel()

	cm := machine.ConfigMap{
		"bar":         machine.Config{Machine: machine.Machine{Cores: 1}},
		"foo.bar":     machine.Config{Machine: machine.Machine{Cores: 2}},
		"foo.bar.baz": machine.Config{Machine: machine.Machine{Cores: 3}},
		"foo.baz":     machine.Config{Machine: machine.Machine{Cores: 4}},
	}

	var mk mocks.Observer
	mk.Test(t)
	mk.On("OnMachines", mock.MatchedBy(func(m machine.Message) bool {
		return m.Kind == machine.MessageStart && m.Index == len(cm)
	})).Return().Once()
	for n := range cm {
		n := n
		mk.On("OnMachines", mock.MatchedBy(func(m machine.Message) bool {
			return m.Kind == machine.MessageRecord && m.Machine.ID.String() == n
		})).Return().Once()
	}
	mk.On("OnMachines", mock.MatchedBy(func(m machine.Message) bool {
		return m.Kind == machine.MessageFinish
	})).Return().Once()

	err := cm.ObserveOn(&mk)
	require.NoError(t, err)

	mk.AssertExpectations(t)
}
