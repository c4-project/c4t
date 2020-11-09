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
	"io/ioutil"
	"math"
	"path"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/filekind"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"
)

// Interpreter is an interpreter for compile recipes.
type Interpreter struct {
	driver Driver
	job    compile.Recipe

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
	fileStack []string
}

var (
	// ErrCompilerConfigNil occurs if a processor is supplied a nil compiler config.
	ErrCompilerConfigNil = errors.New("compiler config nil")
	// ErrBadOp occurs if a processor is supplied an unknown opcode.
	ErrBadOp = errors.New("bad opcode")
	// ErrFileUnavailable occurs if an instruction specifies a file that has been consumed, or wasn't available.
	ErrFileUnavailable = errors.New("file not available")
	// ErrObjOverflow occurs if too many object files are created.
	ErrObjOverflow = errors.New("object file count overflow")
)

// NewInterpreter creates a new recipe processor using the compiler driver d and job j.
func NewInterpreter(d Driver, j compile.Recipe, os ...IOption) (*Interpreter, error) {
	if d == nil {
		return nil, ErrDriverNil
	}
	if j.Compiler == nil {
		return nil, ErrCompilerConfigNil
	}

	p := Interpreter{driver: d, job: j, logw: ioutil.Discard, maxobjs: math.MaxUint64}
	IOptions(os...)(&p)

	p.inPool = initPool(p.job.In)
	// Assuming that the usual case is that every file in the pool gets put in the stack.
	p.fileStack = make([]string, 0, len(p.inPool))

	return &p, nil
}

// IOption is the type of options to the interpreter.
type IOption func(*Interpreter)

// IOptions bundles the options os into one option.
func IOptions(os ...IOption) IOption {
	return func(i *Interpreter) {
		for _, o := range os {
			o(i)
		}
	}
}

// ILogTo logs compiler output to w.
func ILogTo(w io.Writer) IOption {
	return func(i *Interpreter) { i.logw = iohelp.EnsureWriter(w) }
}

// SetMaxObjs sets the maximum number of object files the interpreter can create.
func SetMaxObjs(cap uint64) IOption {
	return func(i *Interpreter) { i.maxobjs = cap }
}

// Interpret processes this processor's compilation recipe using ctx for timeout and cancellation.
// It resumes from the last position where interpretation halted.
func (p *Interpreter) Interpret(ctx context.Context) error {
	ninst := len(p.job.Recipe.Instructions)
	for p.pc < ninst {
		if err := p.processInstruction(ctx, p.job.Recipe.Instructions[p.pc]); err != nil {
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
	p.fileStack = append(p.fileStack, file)
}

func (p *Interpreter) compileObj(ctx context.Context, npops int) error {
	n, err := p.freshObj()
	if err != nil {
		return err
	}
	if err := p.compile(ctx, n, compile.Obj, npops); err != nil {
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
	return path.Join(p.job.Recipe.Dir, file), nil
}

func (p *Interpreter) compileExe(ctx context.Context, npops int) error {
	return p.compile(ctx, p.job.Out, compile.Exe, npops)
	// We don't push the binary onto the file stack.
}

func (p *Interpreter) compile(ctx context.Context, out string, kind compile.Kind, npops int) error {
	if err := p.driver.RunCompiler(ctx, p.singleCompile(out, kind, npops), p.logw); err != nil {
		return err
	}
	return nil
}

func (p *Interpreter) singleCompile(out string, kind compile.Kind, npops int) compile.Single {
	return compile.New(p.job.Compiler, out, p.popStack(npops)...).Single(kind)
}

func (p *Interpreter) popStack(npops int) []string {
	lfs := len(p.fileStack)
	if npops <= 0 || lfs < npops {
		npops = lfs
	}
	cut := lfs - npops

	var fs []string
	fs, p.fileStack = p.fileStack[cut:], p.fileStack[:cut]
	return fs
}

// initPool creates a pool with each path in paths set as available.
func initPool(paths []string) map[string]bool {
	pool := make(map[string]bool, len(paths))
	for _, p := range paths {
		pool[p] = true
	}
	return pool
}
