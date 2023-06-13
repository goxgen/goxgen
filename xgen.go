package goxgen

import (
	"context"
	"embed"
	"github.com/goxgen/goxgen/utils"
)

const (
	ContextPrefix       = "XGEN_CONTEXT_"
	GeneratedFilePrefix = "generated_xgen_"
)

// ContextKey is a key for context
type ContextKey string

// GeneratorContextKey is a key for generator context
var GeneratorContextKey = ContextKey(ContextPrefix + "GENERATOR")

//go:embed templates/*
var templatesFS embed.FS

// Xgen main struct
type Xgen struct {
	PackageName *string
	Projects    []Project
	CLI         *CLI
}

// NewXgen creates a new Xgen instance
// it creates a new Xgen instance with projects and CLI
func NewXgen(options ...XgenOption) *Xgen {
	xgen := &Xgen{
		CLI: &CLI{
			OutputDir: StringP("."),
		},
	}

	for _, opt := range options {
		if err := opt(xgen); err != nil {
			panic(err)
		}
	}

	xgen.CLI.ParentPackageName = xgen.PackageName

	return xgen
}

// prepareGeneratorContext prepares generator context
// it is used to pass data between generators
// for example, gqlgen or project generators can use data from xgen
func (x *Xgen) prepareGeneratorContext(ctx context.Context) context.Context {
	genCtx := Context{
		ParentPackageName: *x.PackageName,
	}

	genCtx.Projects = append(genCtx.Projects, x.Projects...)

	return context.WithValue(ctx, GeneratorContextKey, genCtx)
}

// Generate generates code
// it is a main function of Xgen
// it generates code for projects and CLI
func (x *Xgen) Generate(ctx context.Context) (err error) {

	err = utils.RemoveFromDirByPatterns("./*/" + GeneratedFilePrefix + "*")
	if err != nil {
		return err
	}

	ctx = x.prepareGeneratorContext(ctx)

	err = NewProjectGenerator(x.Projects...).Generate(ctx)
	if err != nil {
		return err
	}

	err = x.CLI.Generate(ctx)
	if err != nil {
		return err
	}

	// exec command `go fmt`
	err = utils.ExecCommand("./", "go", "mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}
