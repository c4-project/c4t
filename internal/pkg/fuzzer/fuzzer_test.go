package fuzzer_test

import (
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
)

// MockPathset mocks the SubjectPather interface.
type MockPathset struct {
	HasPrepared   bool
	SubjectCycles []fuzzer.SubjectCycle
}

func (m *MockPathset) Prepare() error {
	m.HasPrepared = true
	return nil
}

func (m *MockPathset) SubjectPaths(sc fuzzer.SubjectCycle) fuzzer.SubjectPathset {
	m.SubjectCycles = append(m.SubjectCycles, sc)
	return fuzzer.SubjectPathset{
		FileLitmus: path.Join("litmus", sc.String()),
		FileTrace:  path.Join("trace", sc.String()),
	}
}

// makeConfig makes a 'valid' fuzzer config.
func makeConfig() fuzzer.Config {
	return fuzzer.Config{
		Driver:        nil,
		Paths:         &MockPathset{},
		CorpusSize:    0,
		SubjectCycles: 0,
		FuzzWorkers:   0,
	}
}

// TestNewFuzzer_PlanNil makes sure fuzzer creation on a nil plan fails.
func TestNewFuzzer_PlanNil(t *testing.T) {
	c := makeConfig()
	_, err := fuzzer.New(&c, nil)
	testhelp.ExpectErrorIs(t, err, plan.ErrNil, "fuzzer.New on nil plan")
}
