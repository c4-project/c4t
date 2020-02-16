package compiler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// compileJob represents the state of a compiler run.
type compileJob struct {
	// MachineID is the ID of the machine.
	MachineID model.ID

	// Compiler points to the compiler to run.
	Compiler *model.NamedCompiler

	// Pathset is the pathset to use for this compiler run.
	Pathset SubjectPather

	// Compiler is the compiler runner to use for running compilers.
	Runner SingleRunner

	// ResCh is the channel to which the compile run should send results.
	ResCh chan<- result

	// Corpus is the corpus to compile.
	Corpus subject.Corpus
}

func (j *compileJob) Compile(ctx context.Context) error {
	if j.Pathset == nil {
		return fmt.Errorf("in job: %w", iohelp.ErrPathsetNil)
	}

	for i := range j.Corpus {
		if err := j.compileSubject(ctx, &j.Corpus[i]); err != nil {
			return err
		}
	}
	return nil
}

func (j *compileJob) compileSubject(ctx context.Context, s *subject.Subject) error {
	h, herr := s.Harness(j.qualifiedArch())
	if herr != nil {
		return herr
	}

	sp := j.Pathset.SubjectPaths(SubjectCompile{CompilerID: j.Compiler.ID, Name: s.Name})

	logf, err := os.Create(sp.Log)
	if err != nil {
		return err
	}

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	cerr := j.Runner.RunCompiler(j.Compiler, h.Paths(), sp.Bin, logf)

	res, rerr := j.makeCompileResult(sp, s, cerr)
	if rerr != nil {
		_ = logf.Close()
		return rerr
	}

	if err := j.sendResult(ctx, res); err != nil {
		_ = logf.Close()
		return err
	}

	return logf.Close()
}

func (j *compileJob) qualifiedArch() model.MachQualID {
	return model.MachQualID{
		MachineID: j.MachineID,
		ID:        j.Compiler.Arch,
	}
}

func (j *compileJob) makeCompileResult(sp subject.CompileFileset, s *subject.Subject, cerr error) (result, error) {
	cr := result{
		CompilerID: j.qualifiedCompiler(),
		Subject:    s,
		CompileResult: subject.CompileResult{
			// Potentially overridden further down.
			Success: false,
			Files:   sp,
		},
	}

	if cerr != nil {
		// If the error was the compiler process failing to run, then we should report that and carry on.
		var perr *exec.ExitError
		if !errors.As(cerr, &perr) {
			return cr, cerr
		}

		return cr, nil
	}

	cr.Success = true
	return cr, nil
}

func (j *compileJob) qualifiedCompiler() model.MachQualID {
	return model.MachQualID{
		MachineID: j.MachineID,
		ID:        j.Compiler.ID,
	}
}

// sendResult tries to send a compile job result to the result channel.
// If the context ctx has been cancelled, it will fail and instead terminate the job.
func (j *compileJob) sendResult(ctx context.Context, result result) error {
	select {
	case j.ResCh <- result:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
