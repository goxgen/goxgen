package plugins

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/projects"
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

type Plugin interface {
	Generate(ctx context.Context) error
	Name() string
}

func PrepareContext(ctx context.Context, cliCtx *Context) context.Context {
	return context.WithValue(ctx, CLIContextKey, cliCtx)
}

// GetContext gets generator context from context
func GetContext(ctx context.Context) (*Context, error) {
	gCtx, ok := ctx.Value(CLIContextKey).(*Context)

	if !ok {
		return nil, fmt.Errorf("failed to CLI generator context")
	}
	return gCtx, nil
}
