// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

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

	assertPlansSimilar(t, p, &p2)
}

// TestPlan_WriteFile_roundTrip exercises WriteFile by doing a round trip into a temporary file, then checking if
// the reconstituted plan is similar.
func TestPlan_WriteFile_roundTrip(t *testing.T) {
	// Is it safe to call t.Parallel on these?
	dir, err := ioutil.TempDir("", "roundTrip")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	tmpfn := filepath.Join(dir, "plan.json")
	p1 := plan.Mock()

	require.NoErrorf(t, p1.WriteFile(tmpfn), "writing to temp plan file %q", tmpfn)

	var p2 plan.Plan
	require.NoErrorf(t, plan.ReadFile(tmpfn, &p2), "reading from temp plan file %q", tmpfn)

	assertPlansSimilar(t, p1, &p2)
}

func assertPlansSimilar(t *testing.T, p1, p2 *plan.Plan) {
	t.Helper()
	// TODO(@MattWindsor91): more comparisons?
	assert.Truef(t,
		p1.Metadata.Creation.Equal(p2.Metadata.Creation),
		"date not equal after round-trip: send=%v, recv=%v", p1.Metadata.Creation, p2.Metadata.Creation)
	assert.Truef(t,
		p1.Machine.ID.Equal(p2.Machine.ID),
		"machine IDs not equal after round-trip: send=%v, recv=%v", p1.Machine.ID, p2.Machine.ID)
}
