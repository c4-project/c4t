package fuzzer

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/sirupsen/logrus"
)

// prepare does various pre-fuzzing checks and preparation steps.
func (f *Fuzzer) prepare(p *plan.Plan) error {
	if p == nil {
		return plan.ErrNil
	}
	f.Plan = *p

	logrus.Infoln("checking viability")
	if err := f.checkViability(); err != nil {
		return err
	}

	logrus.Infoln("preparing directories")
	return iohelp.Mkdirs(f.Paths)
}

// checkViability does some pre-flight checks.
func (f *Fuzzer) checkViability() error {
	if f.Paths == nil {
		return iohelp.ErrPathsetNil
	}

	if f.SubjectCycles <= 0 {
		return fmt.Errorf("%w: non-positive subject cycle amount", model.ErrSmallCorpus)
	}

	nsubjects, nruns := f.count()
	if nsubjects <= 0 {
		return model.ErrNoCorpus
	}

	// Note that this inequality 'does the right thing' when f.CorpusSize = 0, ie no corpus size requirement.
	if nruns < f.CorpusSize {
		return fmt.Errorf("%w: projected corpus size %d, want %d", model.ErrSmallCorpus, nruns, f.CorpusSize)
	}

	return nil
}
