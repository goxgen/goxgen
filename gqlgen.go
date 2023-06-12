package goxgen

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/gqlgen_plugins"
	"github.com/goxgen/goxgen/utils"
	"path"
)

var GqlgenContextKey = ContextKey(ContextPrefix + "GQLGEN")

type GqlgenContext struct {
	ConfigOverrideCallback func(cfg *config.Config) error
}

func GetGqlgenContext(ctx context.Context) *GqlgenContext {
	if ctx.Value(GqlgenContextKey) != nil {
		return ctx.Value(GqlgenContextKey).(*GqlgenContext)
	}
	return nil
}

func NewGqlgenContext(ctx context.Context, gqlgenContext GqlgenContext) context.Context {
	return context.WithValue(ctx, GqlgenContextKey, &gqlgenContext)
}

func GenerateProjectGqlgenSet(ctx context.Context, project Project) error {

	xgenContext, err := GetXgenContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get xgen context: %w", err)
	}
	gqlgenCtx := GetGqlgenContext(ctx)

	cfg := config.DefaultConfig()

	outputDir := PString(project.GetOutputDir())
	packageName := PString(project.GetName())

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

	cfg.AutoBind = append(cfg.AutoBind, "github.com/goxgen/goxgen/gqlgen_bind")

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
