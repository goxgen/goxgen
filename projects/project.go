package projects

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

// ContextKey is a key for context
type ContextKey int

const contextKey ContextKey = iota

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
	gCtx, ok := ctx.Value(contextKey).(*Context)

	if !ok {
		return nil, fmt.Errorf("project context not found in context")
	}
	return gCtx, nil
}

// Project is a project configuration
type Project interface {
	GetType() string
	TestsDirectory() string
	ModelMutationHook(b *modelgen.ModelBuild) *modelgen.ModelBuild
	ConstraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error)
	SchemaHook(schema *ast.Schema) error
	SchemaDocumentHook(schemaDocument *ast.SchemaDocument) error
	ConfigOverride(cfg *config.Config) error
	Init(Name string, ParentPackageName string, GeneratedFilePrefix string) error
}

// RunProjectGoGenCommand runs the go generate command for the project.
func RunProjectGoGenCommand(dir string) error {
	outputDir := dir
	return utils.ExecCommand(outputDir, "go", "generate")
}

func PrepareCommonContext(ctx context.Context, projCtx *Context) context.Context {
	return context.WithValue(ctx, contextKey, projCtx)
}
