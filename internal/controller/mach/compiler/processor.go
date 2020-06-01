// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

type Processor struct {
	driver SingleRunner
	job    compile.Recipe

	// logw is the writer used for compiler outputs.
	logw io.Writer
	// inPool maps each input file to a Boolean that is true if it hasn't been consumed yet.
	inPool map[string]bool
	// fileStack is the file stack.
	fileStack []string
}

var (
	// ErrCompilerConfigNil occurs if a processor is supplied a nil compiler config.
	ErrCompilerConfigNil = errors.New("compiler config nil")
	// ErrBadOp occurs if a processor is supplied an unknown opcode.
	ErrBadOp = errors.New("bad opcode")
	// ErrFileUnavailable occurs if an instruction specifies a file that has been consumed, or wasn't available.
	ErrFileUnavailable = errors.New("file not available")
)

// NewProcessor creates a new recipe processor using the compiler driver d and job j.
func NewProcessor(d SingleRunner, j compile.Recipe, logw io.Writer) (*Processor, error) {
	if d == nil {
		return nil, ErrDriverNil
	}
	if j.Compiler == nil {
		return nil, ErrCompilerConfigNil
	}
	p := Processor{driver: d, job: j, logw: logw}
	return &p, nil
}

// Process processes this processor's compilation recipe using ctx for timeout and cancellation.
func (p *Processor) Process(ctx context.Context) error {
	p.inPool = initPool(p.job.In)
	// Assuming that the usual case is that every file in the pool gets put in the stack.
	p.fileStack = make([]string, 0, len(p.inPool))

	for _, i := range p.job.Instructions {
		if err := p.processInstruction(ctx, i); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processInstruction(ctx context.Context, i recipe.Instruction) error {
	switch i.Op {
	case recipe.Nop:
		return nil
	case recipe.PushInput:
		return p.pushInput(i.File)
	case recipe.PushInputs:
		return p.pushInputs(i.FileKind)
	case recipe.CompileObj:
		return p.compileObj()
	case recipe.CompileBin:
		return p.compileBin(ctx)
	default:
		return fmt.Errorf("%w: unknown instruction %s", ErrBadOp, i.Op)
	}
}

func (p *Processor) pushInput(file string) error {
	if !p.inPool[file] {
		return fmt.Errorf("%w: %q", ErrFileUnavailable, file)
	}
	p.pushInputRaw(file)
	return nil
}

func (p *Processor) pushInputs(kind filekind.Kind) error {
	for file, ok := range p.inPool {
		if ok && filekind.GuessFromFile(file).Matches(kind) {
			p.pushInputRaw(file)
		}
	}
	return nil
}

func (p *Processor) pushInputRaw(file string) {
	p.inPool[file] = false
	p.fileStack = append(p.fileStack, file)
}

func (p *Processor) compileObj() error {
	// TODO(@MattWindsor91): implement this
	return errors.New("compile to obj not yet implemented")
}

func (p *Processor) compileBin(ctx context.Context) error {
	// TODO(@MattWindsor91): split these two different jobs up
	return p.driver.RunCompiler(ctx, compile.Single{
		Compile: compile.Compile{
			Compiler: p.job.Compiler,
			In:       p.fileStack,
			Out:      p.job.Out,
		},
		Kind: compile.Exe,
	}, p.logw)
}

// initPool creates a pool with each path in paths set as available.
func initPool(paths []string) map[string]bool {
	pool := make(map[string]bool, len(paths))
	for _, p := range paths {
		pool[p] = true
	}
	return pool
}

type ProcessorOption func(*Processor) error
