package compiler

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"golang.org/x/sync/errgroup"
)

func (r *Compiler) compile(ctx context.Context) error {
	eg, ectx := errgroup.WithContext(ctx)

	resCh := make(chan Result)

	for ids, c := range r.Plan.Compilers {
		cr := compileJob{
			CompilerID: model.IDFromString(ids),
			Compiler:   c,
			Runner:     r.Runner,
			ResCh:      resCh,
		}
		eg.Go(func() error {
			return cr.Compile(ectx)
		})
	}

	eg.Go(func() error { return handleCompileResults(ectx) })
	return eg.Wait()
}

// compileJob represents the state of a compiler run.
type compileJob struct {
	// CompilerID is the ID of the compiler.
	CompilerID model.ID

	// Compiler points to the compiler to run.
	Compiler model.Compiler

	// Pathset is the pathset to use for this compiler run.
	Pathset *Pathset

	// Compiler is the compiler runner to use for running compilers.
	Runner CompilerRunner

	// ResCh is the channel to which the compile run should send results.
	ResCh chan<- Result

	// Corpus is the corpus to compile.
	Corpus model.Corpus
}

func (j *compileJob) Compile(ctx context.Context) error {
	for _, s := range j.Corpus {
		if err := j.compileSubject(ctx, &s); err != nil {
			return err
		}
	}
	return nil
}

func (j *compileJob) compileSubject(ctx context.Context, s *model.Subject) error {
	h, herr := s.Harness(j.CompilerID, j.Compiler.Arch)
	if herr != nil {
		return herr
	}

	binp, logp := j.Pathset.OnCompiler(j.CompilerID, s.Name)

	logf, err := os.Create(logp)
	if err != nil {
		return err
	}

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	cerr := j.Runner.RunCompiler(j.CompilerID, h.Paths(), binp, logf)

	res, rerr := j.makeCompileResult(binp, logp, cerr)
	if rerr != nil {
		_ = logf.Close()
		return err
	}

	if err := j.sendResult(ctx, res); err != nil {
		_ = logf.Close()
		return err
	}

	return logf.Close()
}

func (j *compileJob) makeCompileResult(binp, logp string, cerr error) (Result, error) {
	cr := Result{
		CompilerID: j.CompilerID,
		Success:    false,
		PathBin:    binp,
		PathLog:    logp,
	}

	if cerr != nil {
		// If the error was the compiler process failing to run, then we should report that and carry on.
		var perr exec.ExitError
		if !errors.Is(cerr, &perr) {
			return cr, cerr
		}

		return cr, nil
	}

	cr.Success = true
	return cr, nil
}

// sendResult tries to send a compile job result to the result channel.
// If the context ctx has been cancelled, it will fail and instead terminate the job.
func (j *compileJob) sendResult(ctx context.Context, result Result) error {
	select {
	case j.ResCh <- result:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func handleCompileResults(_ context.Context) error {
	return nil
}
