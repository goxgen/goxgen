package goxgen

import (
	"context"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"os"
	"path"
)

// Project is a project configuration
type Project interface {
	GetName() *string
	GetOutputDir() *string
	GetGraphqlURL() *string
	HandleGeneration(ctx context.Context, data *ProjectGeneratorData) error
}

type ProjectWithCustomTemplateData interface {
	Project
	PrepareCustomTemplateData(ctx context.Context, data *ProjectGeneratorData) error
}

type ProjectWithGqlgen interface {
	Project
	GqlgenProject() bool
}

// ProjectGeneratorData is a template data for project
type ProjectGeneratorData struct {
	Name              string
	ParentPackageName string
}

func GenerateProject(ctx context.Context, project Project) error {
	gCtx, err := GetGeneratorContext(ctx)
	if err != nil {
		return err
	}

	templateData := &ProjectGeneratorData{
		ParentPackageName: gCtx.ParentPackageName,
		Name:              PString(project.GetName()),
	}

	if projectWithTemplateData, ok := project.(ProjectWithCustomTemplateData); ok {
		err := projectWithTemplateData.PrepareCustomTemplateData(ctx, templateData)
		if err != nil {
			return err
		}
	}

	err = project.HandleGeneration(ctx, templateData)
	if err != nil {
		return err
	}

	return nil
}

func GenerateProjectGqlgenSet(ctx context.Context, project Project, customConfigCallback func(cfg *config.Config) error) error {
	if projectWithGqlgen, ok := project.(ProjectWithGqlgen); ok && projectWithGqlgen.GqlgenProject() {

		cfg := config.DefaultConfig()

		outputDir := PString(project.GetOutputDir())
		packageName := PString(project.GetName())

		cfg.SchemaFilename = config.StringList{
			path.Join(outputDir, "*.graphql"),
			path.Join(outputDir, "*.graphqls"),
		}

		cfg.Exec.Filename = path.Join(outputDir, "generated_gqlgen.go")
		cfg.Exec.DirName = outputDir

		cfg.Resolver.DirName = outputDir
		cfg.Resolver.FilenameTemplate = "{name}.graphql.resolver.go"

		// mark standard resolver.go as redundant to be deleted after gqlgen generate
		// because we don't need it anymore, we create new Resolver instance in `project_handlers.go`
		cfg.Resolver.Filename = path.Join(outputDir, "_redundant_resolver")
		cfg.Resolver.Package = packageName
		cfg.Resolver.Layout = "follow-schema"

		cfg.Model.Package = packageName
		cfg.Model.Filename = path.Join(outputDir, "generated_gqlgen_models.go")

		err := customConfigCallback(cfg)
		if err != nil {
			return err
		}

		err = config.CompleteConfig(cfg)
		if err != nil {
			return err
		}

		if err = api.Generate(cfg); err != nil {
			return err
		}

		// delete redundant _redundant_resolver
		redundantResolverFile := path.Join(PString(project.GetOutputDir()), "_redundant_resolver")
		if _, err := os.Stat(redundantResolverFile); err == nil {
			err = os.Remove(redundantResolverFile)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func RunProjectGoGenCommand(project Project) error {
	outputDir := PString(project.GetOutputDir())
	return ExecCommand(outputDir, "go", "generate")
}
