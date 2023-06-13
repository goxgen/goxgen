package goxgen

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/directives"
	"github.com/goxgen/goxgen/gqlgen_plugins"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"path"
)

var GraphqlContextKey = ContextKey(ContextPrefix + "GQLGEN")

type GraphqlContext struct {
	ConfigOverrideCallback func(cfg *config.Config) error
}

// GetGraphqlContext returns the graphql context from the context.
func GetGraphqlContext(ctx context.Context) *GraphqlContext {
	if ctx.Value(GraphqlContextKey) != nil {
		return ctx.Value(GraphqlContextKey).(*GraphqlContext)
	}
	return nil
}

// NewGraphqlContext returns a new context with the graphql context.
func NewGraphqlContext(ctx context.Context, gqlgenContext GraphqlContext) context.Context {
	return context.WithValue(ctx, GraphqlContextKey, &gqlgenContext)
}

// generateDirectivesSet generates a graphql file with all the Xgen directives.
func generateDirectivesSet(outputDir string) error {
	schemaGenerator := graphql.SchemaGenerator{
		Path: path.Join(outputDir, GeneratedFilePrefix+"directives.graphql"),
		SchemaHooks: []graphql.SchemaHook{
			func(schema *ast.Schema) error {
				for _, directive := range directives.All {
					schema.Directives[directive.Name] = directive
				}
				return nil
			},
		},
	}
	return schemaGenerator.GenerateOutput()
}

// GenerateProjectGraphqlSet generates a graphql set for the project.
// Using the gqlgen library.
func GenerateProjectGraphqlSet(ctx context.Context, project Project) error {
	outputDir := PString(project.GetOutputDir())
	packageName := PString(project.GetName())

	xgenContext, err := GetContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get xgen context: %w", err)
	}

	gqlgenCtx := GetGraphqlContext(ctx)

	err = generateDirectivesSet(outputDir)
	if err != nil {
		return fmt.Errorf("failed to generate directives set: %w", err)
	}

	cfg := config.DefaultConfig()

	cfg.SchemaFilename = config.StringList{
		path.Join(outputDir, "*.graphql"),
		path.Join(outputDir, "*.graphqls"),
	}

	gqlgenPackage := "generated_gqlgen"
	gqlgenPath := path.Join(outputDir, gqlgenPackage)

	err = utils.RemoveFromDirByPatterns(gqlgenPath)
	if err != nil {
		return fmt.Errorf("failed to remove old gqlgen_generated files: %w", err)
	}

	cfg.Exec.Package = gqlgenPackage
	cfg.Exec.Filename = path.Join(gqlgenPath, "generated_gqlgen.go")

	cfg.Resolver.DirName = outputDir
	cfg.Resolver.FilenameTemplate = "{name}.graphql.resolver.go"

	// mark standard resolver.go as redundant to be deleted after gqlgen generate
	// because we don't need it anymore, we create new Resolver instance in `project_handlers.go`
	cfg.Resolver.Filename = path.Join(outputDir, "_redundant_resolver")
	cfg.Resolver.Package = packageName
	cfg.Resolver.Layout = "follow-schema"

	cfg.Model.Package = gqlgenPackage
	cfg.Model.Filename = path.Join(gqlgenPath, "generated_gqlgen_models.go")

	//cfg.AutoBind = append(cfg.AutoBind, "github.com/goxgen/goxgen/gqlgen_bind")

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

	injectorPlugin, defTypesInjectorPlugin := gqlgen_plugins.NewPlugin(
		outputDir,
		packageName,
		xgenContext.ParentPackageName+"/"+PString(project.GetName()),
		GeneratedFilePrefix,
	)
	if err = api.Generate(cfg, api.AddPlugin(injectorPlugin)); err != nil {
		return fmt.Errorf("failed to generate gqlgen files: %w", err)
	}

	cfg.Model.Filename = path.Join(outputDir, cfg.Exec.Package, "generated_gqlgen_defs_models.go")
	if err = api.Generate(cfg, api.AddPlugin(defTypesInjectorPlugin)); err != nil {
		return fmt.Errorf("failed to generate gqlgen files phase 2: %w", err)
	}

	// delete redundant _redundant_resolver
	err = utils.RemoveFromDirByPatterns(path.Join(PString(project.GetOutputDir()), "_redundant_resolver"))
	if err != nil {
		return fmt.Errorf("failed to remove _redundant_resolver: %w", err)
	}

	return nil
}
