// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// TestPlan_Write_roundTrip exercises Write by doing a round-trip and checking if the reconstituted plan is similar.
func TestPlan_Write_roundTrip(t *testing.T) {
	t.Parallel()

	p := plan.Mock()

	var b bytes.Buffer
	if err := p.Write(&b); err != nil {
		t.Fatal("error dumping:", err)
	}

	var p2 plan.Plan
	if err := plan.Read(&b, &p2); err != nil {
		t.Fatal("error un-dumping:", err)
	}

	// TODO(@MattWindsor91): more comparisons?
	assert.Truef(t,
		p.Metadata.Creation.Equal(p2.Metadata.Creation),
		"date not equal after round-trip: send=%v, recv=%v", p.Metadata.Creation, p2.Metadata.Creation)
}
