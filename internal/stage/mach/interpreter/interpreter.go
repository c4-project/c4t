// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package interpreter contains the recipe interpreter for the machine node.
package interpreter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path"
	"path/filepath"

	"github.com/MattWindsor91/c4t/internal/model/service/compiler"

	"github.com/MattWindsor91/c4t/internal/model/filekind"

	"github.com/MattWindsor91/c4t/internal/model/recipe"
)

// ErrDriverNil occurs when the compiler tries to use the nil pointer as its single-compile driver.
var ErrDriverNil = errors.New("driver nil")

// Driver is the interface of things that can run compilers.
type Driver interface {
	// RunCompiler runs the compiler job j.
	// If applicable, errw will be connected to the compiler's standard error.
	//
	// Implementors should note that the paths in j are slash-paths, and will need converting to filepaths.
	RunCompiler(ctx context.Context, j compiler.Job, errw io.Writer) error
}

//go:generate mockery --name=Driver

// Interpreter is an interpreter for compile recipes.
type Interpreter struct {
	driver Driver
	// compiler is the compiler configuration.
	compiler *compiler.Configuration
	// ofile is the output filepath.
	ofile string
	// recipe is the recipe to interpret.
	recipe recipe.Recipe

	// pc is the program counter.
	pc int
	// nobjs is the number of object files created so far by the processor.
	nobjs uint64
	// maxobjs is the maximum permitted number of object files.
	maxobjs uint64
	// logw is the writer used for compiler outputs.
	logw io.Writer
	// inPool maps each input file to a Boolean that is true if it hasn't been consumed yet.
	inPool map[string]bool
	// fileStack is the file stack.
	fileStack stack
}

var (
	// ErrCompilerConfigNil occurs if an interpreter is supplied a nil compiler config.
	ErrCompilerConfigNil = errors.New("compiler config nil")
	// ErrBadOp occurs if an interpreter is supplied an unknown opcode.
	ErrBadOp = errors.New("bad opcode")
	// ErrBadOutput occurs if an interpreter is asked to output something that isn't compatible with its output spec.
	ErrBadOutput = errors.New("bad output type")
	// ErrFileUnavailable occurs if an instruction specifies a file that has been consumed, or wasn't available.
	ErrFileUnavailable = errors.New("file not available")
	// ErrObjOverflow occurs if too many object files are created.
	ErrObjOverflow = errors.New("object file count overflow")
)

// NewInterpreter creates a new recipe processor using the compiler driver d, configuration c, runner r, and job j.
func NewInterpreter(d Driver, c *compiler.Configuration, ofile string, r recipe.Recipe, os ...Option) (*Interpreter, error) {
	if d == nil {
		return nil, ErrDriverNil
	}
	if c == nil {
		return nil, ErrCompilerConfigNil
	}

	p := Interpreter{driver: d, compiler: c, ofile: ofile, recipe: r, logw: ioutil.Discard, maxobjs: math.MaxUint64}
	Options(os...)(&p)

	p.inPool = p.initPool()
	// Assuming that the usual case is that every file in the pool gets put in the stack.
	p.fileStack = make([]string, 0, len(p.inPool))

	return &p, nil
}

// Interpret processes this processor's compilation recipe using ctx for timeout and cancellation.
// It resumes from the last position where interpretation halted.
func (p *Interpreter) Interpret(ctx context.Context) error {
	ninst := len(p.recipe.Instructions)
	for p.pc < ninst {
		if err := p.processInstruction(ctx, p.recipe.Instructions[p.pc]); err != nil {
			return err
		}
		p.pc++
	}
	return nil
}

func (p *Interpreter) processInstruction(ctx context.Context, i recipe.Instruction) error {
	switch i.Op {
	case recipe.Nop:
		return nil
	case recipe.PushInput:
		return p.pushInput(i.File)
	case recipe.PushInputs:
		return p.pushInputs(i.FileKind)
	case recipe.CompileObj:
		return p.compileObj(ctx, i.NPops)
	case recipe.CompileExe:
		return p.compileExe(ctx, i.NPops)
	default:
		return fmt.Errorf("%w: unknown instruction %s", ErrBadOp, i.Op)
	}
}

func (p *Interpreter) pushInput(file string) error {
	if !p.inPool[file] {
		return fmt.Errorf("%w: %q", ErrFileUnavailable, file)
	}
	p.pushInputRaw(file)
	return nil
}

func (p *Interpreter) pushInputs(kind filekind.Kind) error {
	for file, ok := range p.inPool {
		if ok && filekind.GuessFromFile(file).Matches(kind) {
			p.pushInputRaw(file)
		}
	}
	return nil
}

func (p *Interpreter) pushInputRaw(file string) {
	p.inPool[file] = false
	p.fileStack.push(file)
}

func (p *Interpreter) compileObj(ctx context.Context, npops int) error {
	n, err := p.freshObj()
	if err != nil {
		return err
	}
	if err := p.compile(ctx, n, compiler.Obj, npops); err != nil {
		return err
	}
	p.fileStack = append(p.fileStack, n)
	return nil
}

func (p *Interpreter) freshObj() (string, error) {
	if p.nobjs == p.maxobjs {
		return "", ErrObjOverflow
	}
	// TODO(@MattWindsor91): filepath?
	file := fmt.Sprintf("obj_%d.o", p.nobjs)
	p.nobjs++
	return path.Join(p.recipe.Dir, file), nil
}

func (p *Interpreter) compileExe(ctx context.Context, npops int) error {
	if p.recipe.Output != recipe.OutExe {
		return fmt.Errorf("%w: cannot compile exe when targeting %q", ErrBadOutput, p.recipe.Output)
	}
	return p.compile(ctx, p.ofile, compiler.Exe, npops)
	// We don't push the binary onto the file stack.
}

func (p *Interpreter) compile(ctx context.Context, out string, kind compiler.Target, npops int) error {
	return p.driver.RunCompiler(ctx, *p.singleCompile(out, kind, npops), p.logw)
}

func (p *Interpreter) singleCompile(out string, kind compiler.Target, npops int) *compiler.Job {
	return compiler.NewJob(kind, p.compiler, out, p.fileStack.pop(npops)...)
}

// initPool creates a pool with each path in paths set as available.
func (p *Interpreter) initPool() map[string]bool {
	pool := make(map[string]bool, len(p.recipe.Files))
	dir := filepath.Clean(p.recipe.Dir)
	for _, file := range p.recipe.Files {
		pool[filepath.Join(dir, file)] = true
	}
	return pool
}
