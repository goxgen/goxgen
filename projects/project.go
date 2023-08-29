package projects

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/utils"
)

// ContextKey is a key for context
type ContextKey string

var ProjectsContextKey = ContextKey("PROJECTS")

// Context is a context for Xgen
// it is used to pass data between generators
type Context struct {
	ParentPackageName   string
	GeneratedFilePrefix string
	Projects            map[string]Project
}

// GetContext gets generator context from context
// it is used to pass data between generators
func GetContext(ctx context.Context) (*Context, error) {
	gCtx, ok := ctx.Value(ProjectsContextKey).(*Context)

	if !ok {
		return nil, fmt.Errorf("project context not found in context")
	}
	return gCtx, nil
}

// Project is a project configuration
type Project interface {
	GetType() *string
	PrepareGraphqlGenerationContext(projCtx *Context, data *ProjectGeneratorData) (context.Context, error)
}

// ProjectWithCustomTemplateData is a project configuration with custom template data
type ProjectWithCustomTemplateData interface {
	Project
	PrepareCustomTemplateData(projCtx *Context, data *ProjectGeneratorData) error
}

// RunProjectGoGenCommand runs the go generate command for the project.
func RunProjectGoGenCommand(dir string) error {
	outputDir := dir
	return utils.ExecCommand(outputDir, "go", "generate")
}

func PrepareCommonContext(ctx context.Context, projCtx *Context) context.Context {
	return context.WithValue(ctx, ProjectsContextKey, projCtx)
}
