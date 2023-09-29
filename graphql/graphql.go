package graphql

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/db"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/graphql/pagination"
	"github.com/goxgen/goxgen/graphql/resource"
	"github.com/goxgen/goxgen/graphql/sort"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"path"
)

// ContextKey is a key for context
type ContextKey string

var GraphqlContextKey = ContextKey("GRAPHQL_CONTEXT")

type GraphqlContext struct {
	ParentPackageName           string
	GeneratedFilePrefix         string
	ConfigOverrideCallback      func(cfg *config.Config) error
	CustomDirectivesDefinitions []*ast.DirectiveDefinition
	CustomSchemaFiles           config.StringList
	SchemaInjectorHooks         []InjectorHook
	GeneratorApiOptions         []api.Option
}

var (
	MainDirectiveDefinitionBundle = &directives.DirectiveDefinitionBundle{
		Object: []*directives.ObjectDirectiveDefinition{
			&resource.Directive,
		},
		InputObject: []*directives.InputObjectDirectiveDefinition{
			&resource.ActionDirective,
			&resource.ListActionDirective,
			&common.ExcludeArgumentFromTypeDirective,
			&common.ToObjectType,
		},
		Field: []*directives.FieldDirectiveDefinition{
			&resource.FieldDirective,
		},
		InputField: []*directives.InputFieldDirectiveDefinition{
			&resource.ActionFieldDirective,
		},
	}
)

// GetGraphqlContext returns the graphql context from the context.
func GetGraphqlContext(ctx context.Context) (*GraphqlContext, error) {
	if ctx.Value(GraphqlContextKey) != nil {
		return ctx.Value(GraphqlContextKey).(*GraphqlContext), nil
	}
	return nil, fmt.Errorf("failed to get graphql context")
}

// PrepareContext returns a new context with the graphql context.
func PrepareContext(ctx context.Context, gqlgenContext GraphqlContext) context.Context {
	return context.WithValue(ctx, GraphqlContextKey, &gqlgenContext)
}

// generateDirectivesSet generates a graphql file with all the Xgen directives.
func generateDirectivesSet(outputDir string, generatedFilePrefix string) error {
	schemaGenerator := generator.NewSchemaGenerator().
		WithPath(path.Join(outputDir, generatedFilePrefix+"directives.graphql")).
		WithSchemaHooks(func(_document *ast.SchemaDocument) error {
			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, resource.AllDefinitions...)
			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, sort.AllDefinitions...)
			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, pagination.AllDefinitions...)
			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, db.AllDefinitions...)
			_document.Directives = append(_document.Directives, MainDirectiveDefinitionBundle.DirectiveDefinitionList()...)
			return nil
		})
	return schemaGenerator.GenerateOutput()
}

// Generate generates a graphql set for the project.
// Using the gqlgen library.
func Generate(ctx context.Context, name string) error {
	gqlgenCtx, err := GetGraphqlContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get graphql context: %w", err)
	}

	err = generateDirectivesSet(name, gqlgenCtx.GeneratedFilePrefix)
	if err != nil {
		return fmt.Errorf("failed to generate directives set: %w", err)
	}

	cfg := config.DefaultConfig()

	cfg.SchemaFilename = append(config.StringList{
		path.Join(name, "schema.*.graphql"),
		path.Join(name, "schema.*.graphqls"),
		path.Join(name, gqlgenCtx.GeneratedFilePrefix+"directives.graphql"),
	}, gqlgenCtx.CustomSchemaFiles...)

	gqlgenPath := path.Join(name, consts.GeneratedGqlgenPackageName)

	err = utils.RemoveFromDirByPatterns(gqlgenPath)
	if err != nil {
		return fmt.Errorf("failed to remove old gqlgen_generated files: %w", err)
	}

	cfg.Exec.Package = consts.GeneratedGqlgenPackageName
	cfg.Exec.Filename = path.Join(gqlgenPath, "generated_gqlgen.go")

	cfg.Resolver.DirName = name
	cfg.Resolver.FilenameTemplate = "{name}.resolver.go"

	cfg.Resolver.Package = name
	cfg.Resolver.Layout = "follow-schema"

	cfg.Model.Package = consts.GeneratedGqlgenPackageName
	cfg.Model.Filename = path.Join(gqlgenPath, "generated_gqlgen_models.go")

	if gqlgenCtx != nil && gqlgenCtx.ConfigOverrideCallback != nil {
		err := gqlgenCtx.ConfigOverrideCallback(cfg)
		if err != nil {
			return fmt.Errorf("failed to override gqlgen config: %w", err)
		}
	}

	err = config.CompleteConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to complete gqlgen config: %w", err)
	}

	injectorPlugin := NewPlugin(
		name,
		gqlgenCtx.ParentPackageName,
		gqlgenCtx.GeneratedFilePrefix,
		gqlgenCtx.SchemaInjectorHooks...,
	)

	if err = api.Generate(
		cfg,
		append(
			[]api.Option{api.AddPlugin(injectorPlugin)},
			gqlgenCtx.GeneratorApiOptions...,
		)...,
	); err != nil {
		return fmt.Errorf("failed to generate gqlgen files phase 1: %w", err)
	}

	return nil
}
