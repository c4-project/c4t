package runner

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ObsParser is the interface of things that can parse test outcomes.
type ObsParser interface {
	// ParseObs parses the observation in reader r into o according to the backend configuration in b.
	// The backend described by b must have been used to produce the testcase outputting r.
	ParseObs(ctx context.Context, b model.Backend, r io.Reader, o *model.Obs) error
}

// Config represents the configuration needed to run a Runner.
type Config struct {
	// Parser handles the parsing of observations.
	Parser ObsParser

	// Paths contains the pathset used for this runner's outputs.
	Paths *Pathset

	// MachineID is the ID of the machine whose compilations we are running.
	MachineID model.ID
}
