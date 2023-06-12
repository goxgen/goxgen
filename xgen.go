package goxgen

import (
	"context"
	"embed"
	"fmt"
	"github.com/goxgen/goxgen/utils"
)

const (
	ContextPrefix       = "XGEN_CONTEXT_"
	GeneratedFilePrefix = "generated_xgen_"
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

type XgenContext struct {
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
	genCtx := XgenContext{
		ParentPackageName: *x.PackageName,
	}

	genCtx.Projects = append(genCtx.Projects, x.Projects...)

	return context.WithValue(ctx, GeneratorContextKey, genCtx)
}

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

func GetXgenContext(ctx context.Context) (*XgenContext, error) {
	gCtx, ok := ctx.Value(GeneratorContextKey).(XgenContext)

	if !ok {
		return nil, fmt.Errorf("failed to get generator context")
	}
	return &gCtx, nil
}
