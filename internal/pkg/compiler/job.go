package compiler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

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

	// ResCh is the channel to which the compile run should send compiled subject records.
	ResCh chan<- corpus.BuilderReq

	// Corpus is the corpus to compile.
	Corpus corpus.Corpus
}

func (j *compileJob) Compile(ctx context.Context) error {
	if j.Pathset == nil {
		return fmt.Errorf("in job: %w", iohelp.ErrPathsetNil)
	}

	return j.Corpus.Each(func(s subject.Named) error {
		return j.compileSubject(ctx, &s)
	})
}

func (j *compileJob) compileSubject(ctx context.Context, s *subject.Named) error {
	h, herr := s.Harness(j.qualifiedArch())
	if herr != nil {
		return herr
	}

	sp := j.Pathset.SubjectPaths(SubjectCompile{CompilerID: j.Compiler.ID, Name: s.Name})

	res, rerr := j.runCompiler(ctx, sp, h)
	if rerr != nil {
		return rerr
	}

	return j.sendResult(ctx, s.Name, res)
}

func (j *compileJob) runCompiler(ctx context.Context, sp subject.CompileFileset, h subject.Harness) (subject.CompileResult, error) {
	logf, err := os.Create(sp.Log)
	if err != nil {
		return subject.CompileResult{}, err
	}

	// Some compiler errors are recoverable, so we don't immediately bail on them.
	cerr := j.Runner.RunCompiler(ctx, j.Compiler, h.Paths(), sp.Bin, logf)

	// We could close the log file here, but we want fatal compiler errors to take priority over log file close errors.
	res, rerr := j.makeCompileResult(sp, cerr)
	if rerr != nil {
		_ = logf.Close()
		return subject.CompileResult{}, rerr
	}

	lerr := logf.Close()
	return res, lerr
}

func (j *compileJob) qualifiedArch() model.MachQualID {
	return model.MachQualID{
		MachineID: j.MachineID,
		ID:        j.Compiler.Arch,
	}
}

// makeCompileResult makes a compile result given a possible compile error cerr and fileset sp.
// It fails if the compile error is considered substantially fatal.
func (j *compileJob) makeCompileResult(sp subject.CompileFileset, cerr error) (subject.CompileResult, error) {
	cr := subject.CompileResult{
		// Potentially overridden further down.
		Success: false,
		Files:   sp,
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

// sendResult tries to send a compile job result to the result channel.
// If the context ctx has been cancelled, it will fail and instead terminate the job.
func (j *compileJob) sendResult(ctx context.Context, name string, result subject.CompileResult) error {
	select {
	case j.ResCh <- j.builderReq(name, result):
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// builderReq makes a builder request for adding r to the subject named name.
func (j *compileJob) builderReq(name string, r subject.CompileResult) corpus.BuilderReq {
	return corpus.BuilderReq{
		Name: name,
		Req: corpus.AddCompileReq{
			CompilerID: j.qualifiedCompiler(),
			Result:     r,
		},
	}
}

// qualifiedCompiler gets the machine-qualified ID of this job's compiler.
func (j *compileJob) qualifiedCompiler() model.MachQualID {
	return model.MachQualID{
		MachineID: j.MachineID,
		ID:        j.Compiler.ID,
	}
}
