package cli

import (
	"context"
	"embed"
	"fmt"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/templates_engine"
)

// ContextKey is a key for context
type ContextKey string

var CLIContextKey = ContextKey("CLI")

// Context is a context for CLI
type Context struct {
	ParentPackageName   string
	GeneratedFilePrefix string
	Projects            map[string]projects.Project
	OutputDir           string
}

// GetContext gets generator context from context
func GetContext(ctx context.Context) (*Context, error) {
	gCtx, ok := ctx.Value(CLIContextKey).(*Context)

	if !ok {
		return nil, fmt.Errorf("failed to CLI generator context")
	}
	return gCtx, nil
}

// CLI contains configuration of xgen-cli
type CLI struct {
}

// templateData contains data for xgen-cli templates
type templateData struct {
	ParentPackageName string
	Projects          map[string]projects.Project
}

//go:embed templates/*
var templatesFS embed.FS

// Generate generates code from templates
func (xc *CLI) Generate(ctx context.Context) error {
	genCtx, err := GetContext(ctx)
	if err != nil {
		return err
	}

	data, err := xc.prepareCLITemplateData(genCtx)
	if err != nil {
		return err
	}

	tb := templates_engine.TemplateBundle{
		TemplateDir: "templates",
		FS:          templatesFS,
		OutputFile:  "./" + genCtx.GeneratedFilePrefix + "cli.go",
		Regenerate:  true,
	}

	return tb.Generate(genCtx.OutputDir, data)
}

func (xc *CLI) prepareCLITemplateData(genCtx *Context) (*templateData, error) {
	return &templateData{
		Projects:          genCtx.Projects,
		ParentPackageName: genCtx.ParentPackageName,
	}, nil
}

func PrepareContext(ctx context.Context, cliCtx *Context) context.Context {

	return context.WithValue(ctx, CLIContextKey, cliCtx)
}
