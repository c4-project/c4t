package corpus

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// BuilderReq is the type of requests to a Builder.
type BuilderReq struct {
	// Name is the name of the subject to add or modify
	Name string

	// Req is the request payload, which will be one of the *Req structs.
	Req interface{}
}

// AddReq is a request to add the given subject to the corpus.
type AddReq subject.Subject

// AddCompileReq is a request to add the given compiler result to the named subject.
type AddCompileReq struct {
	// CompilerID is the machine-qualified ID of the compiler that produced this result.
	CompilerID model.MachQualID

	// Result is the compile result.
	Result subject.CompileResult
}
