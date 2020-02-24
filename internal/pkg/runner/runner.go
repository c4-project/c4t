// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package runner contains the part of act-tester that runs compiled harness binaries and interprets their output.
package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

var (
	// ErrNoBin occurs when a successful compile result	has no binary path attached.
	ErrNoBin = errors.New("no binary in compile result")

	// ErrConfigNil occurs when we try to construct a Runner using a nil Config.
	ErrConfigNil = errors.New("config nil")
)

// Runner contains information necessary to run a plan's compiled test cases.
type Runner struct {
	// l is the logger for this runner.
	l *log.Logger

	// plan is the plan on which this runner is operating.
	plan plan.Plan

	// mid is the ID of the machine on which this runner is operating.
	mid model.ID

	// mach is a copy of the specific machine (in plan) on which this runner is operating.
	mach plan.MachinePlan

	// conf is the configuration used to build this runner.
	conf Config
}

// New creates a new batch compiler instance using the config c and plan p.
// It can fail if various safety checks fail on the config,
// or if there is no obvious machine that the compiler can target.
func New(c *Config, p *plan.Plan) (*Runner, error) {
	if c == nil {
		return nil, ErrConfigNil
	}
	if p == nil {
		return nil, plan.ErrNil
	}
	mid, mach, err := p.Machine(c.MachineID)
	if err != nil {
		return nil, err
	}

	r := Runner{
		conf: *c,
		mid:  mid,
		mach: mach,
		plan: *p,
		l:    iohelp.EnsureLog(c.Logger),
	}

	if err := r.check(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Runner) check() error {
	if len(r.plan.Corpus) == 0 {
		return corpus.ErrNoCorpus
	}
	return nil
}

// Run runs the runner.
func (r *Runner) Run(ctx context.Context) (*Result, error) {
	res := Result{
		Time:     time.Now(),
		Subjects: make(map[string]SubjectResult, len(r.plan.Corpus)),
	}

	err := r.plan.Corpus.Map(func(named *subject.Named) error {
		var err error
		res.Subjects[named.Name], err = r.runSubject(ctx, named)
		return err
	})
	return &res, err
}

func (r *Runner) runSubject(ctx context.Context, s *subject.Named) (SubjectResult, error) {
	r.l.Println("running subject:", s.Name)

	var err error
	res := SubjectResult{Compilers: make(map[string]CompilerResult, len(s.Compiles))}

	for cidstr, c := range s.Compiles {
		cid := model.IDFromString(cidstr)

		if res.Compilers[cidstr], err = r.runCompile(ctx, s, cid, &c); err != nil {
			return res, err
		}
	}
	return res, nil
}

func (r *Runner) runCompile(ctx context.Context, s *subject.Named, cid model.ID, c *subject.CompileResult) (CompilerResult, error) {
	if !c.Success {
		return CompilerResult{Status: StatusCompileFail}, nil
	}

	bin := c.Files.Bin
	if bin == "" {
		return CompilerResult{Status: StatusUnknown}, fmt.Errorf("%w: subject=%s, compiler=%s", ErrNoBin, s.Name, cid.String())
	}

	obs, runErr := r.runAndParseBin(ctx, bin)
	status, err := StatusOfObs(obs, runErr)

	return CompilerResult{Status: status, Obs: obs}, err
}

// runAndParseBin runs the binary at bin and parses its result into an observation struct.
func (r *Runner) runAndParseBin(ctx context.Context, bin string) (*model.Obs, error) {
	tctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	cmd := exec.CommandContext(tctx, bin)
	obsr, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var obs model.Obs
	if err := r.conf.Parser.ParseObs(ctx, r.mach.Backend, obsr, &obs); err != nil {
		_ = cmd.Wait()
		return nil, err
	}

	err = cmd.Wait()
	return &obs, err
}
