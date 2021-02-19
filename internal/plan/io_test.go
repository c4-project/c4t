// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/plan"
)

// TestPlan_Write_roundTrip exercises Write by doing a round-trip and checking if the reconstituted plan is similar.
func TestPlan_Write_roundTrip(t *testing.T) {
	t.Parallel()

	cases := map[string]plan.WriteFlag{
		"mach":     plan.WriteNone,
		"mach+gz":  plan.WriteCompress,
		"human":    plan.WriteHuman,
		"human+gz": plan.WriteHuman | plan.WriteCompress,
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			p := plan.Mock()
			var b bytes.Buffer
			if err := p.Write(&b, c); err != nil {
				t.Fatal("error dumping:", err)
			}

			r := bytes.NewReader(b.Bytes())
			var p2 plan.Plan
			if err := plan.ReadMagic(r, &p2); err != nil {
				t.Fatal("error un-dumping:", err)
			}

			assertPlansSimilar(t, p, &p2)
		})
	}
}

// TestPlan_WriteFile_roundTrip exercises WriteFile by doing a round trip into a temporary file, then checking if
// the reconstituted plan is similar.
func TestPlan_WriteFile_roundTrip(t *testing.T) {
	// Is it safe to call t.Parallel on these?
	dir := t.TempDir()
	defer func() { _ = os.RemoveAll(dir) }()

	tmpfn := filepath.Join(dir, "small-ok.json")

	cases := map[string]plan.WriteFlag{
		"mach":     plan.WriteNone,
		"mach+gz":  plan.WriteCompress,
		"human":    plan.WriteHuman,
		"human+gz": plan.WriteHuman | plan.WriteCompress,
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			// See above in re parallel.

			p1 := plan.Mock()

			require.NoErrorf(t, p1.WriteFile(tmpfn, c), "writing to temp plan file %q", tmpfn)

			var p2 plan.Plan
			require.NoErrorf(t, plan.ReadFile(tmpfn, &p2), "reading from temp plan file %q", tmpfn)

			assertPlansSimilar(t, p1, &p2)
		})
	}
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
