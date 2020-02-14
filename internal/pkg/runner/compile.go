package runner

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"golang.org/x/sync/errgroup"
)

func (r *Runner) compile(ctx context.Context) error {
	eg, ectx := errgroup.WithContext(ctx)

	resCh := make(chan compileResult)

	for ids, c := range r.Plan.Compilers {
		cr := compileRun{
			ID:       model.IDFromString(ids),
			Compiler: c,
			Runner:   r.Compiler,
			ResCh:    resCh,
		}
		eg.Go(func() error {
			return cr.Compile(ectx)
		})
	}

	eg.Go(func() error { return handleCompileResults(ectx) })
	return eg.Wait()
}

// compileRun represents the state of a compiler run.
type compileRun struct {
	// ID is the ID of the compiler.
	ID model.ID

	// Compiler points to the compiler to run.
	Compiler model.Compiler

	// Pathset is the pathset to use for this compiler run.
	Pathset *Pathset

	// Runner is the compiler runner to use for running compilers.
	Runner CompilerRunner

	// ResCh is the channel to which the compile run should send results.
	ResCh chan<- compileResult

	// Corpus is the corpus to compile.
	Corpus model.Corpus
}

func (cr *compileRun) Compile(ctx context.Context) error {
	for _, s := range cr.Corpus {
		if err := cr.compileSubject(ctx, &s); err != nil {
			return err
		}
	}
	return nil
}

func (cr *compileRun) compileSubject(ctx context.Context, s *model.Subject) error {
	h, herr := s.Harness(cr.ID, cr.Compiler.Arch)
	if herr != nil {
		return herr
	}

	if err := cr.Runner.RunCompiler(cr.ID, h.Paths(), ""); err != nil {
		return err
	}
	return nil
}

// compileResult represents the output of a compiler goroutine.
type compileResult struct {
	// ID is the ID of the compiler used to compile this binary.
	ID model.ID

	// Success gets whether the compilation succeeded (possibly with errors).
	Success bool

	// PathBin, on success, provides the path to the compiled binary.
	PathBin string

	// PathLog provides the path to the compiler's stderr log.
	PathLog string
}

func handleCompileResults(ctx context.Context) error {
	return nil
}
