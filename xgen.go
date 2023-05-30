package goxgen

import (
	"context"
	"embed"
	"fmt"
)

const (
	ContextPrefix = "XGEN_CONTEXT_"
)

type ContextKey string

var GeneratorContextKey = ContextKey(ContextPrefix + "GENERATOR")

//go:embed templates/*
var templatesFS embed.FS

// Xgen main struct
type Xgen struct {
	PackageName *string
	Projects    []Project
	CLI         *CLI
}

type GeneratorContext struct {
	ParentPackageName string
	Projects          []Project
}

// NewXgen creates a new Xgen instance
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

func (x *Xgen) prepareGeneratorContext(ctx context.Context) context.Context {
	genCtx := GeneratorContext{
		ParentPackageName: *x.PackageName,
	}

	genCtx.Projects = append(genCtx.Projects, x.Projects...)

	return context.WithValue(ctx, GeneratorContextKey, genCtx)
}

func (x *Xgen) Generate(ctx context.Context) error {

	ctx = x.prepareGeneratorContext(ctx)

	for _, p := range x.Projects {
		err := GenerateProject(ctx, p)
		if err != nil {
			return err
		}
	}

	err := x.CLI.Generate(ctx)
	if err != nil {
		return err
	}

	// exec command `go fmt`
	err = ExecCommand("./", "go", "mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}

func GetGeneratorContext(ctx context.Context) (*GeneratorContext, error) {
	gCtx, ok := ctx.Value(GeneratorContextKey).(GeneratorContext)

	if !ok {
		return nil, fmt.Errorf("failed to get generator context")
	}
	return &gCtx, nil
}
