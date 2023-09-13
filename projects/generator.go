package projects

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"path"
)

// ProjectGenerator is a wrapper for project configuration
// That takes context from xgen and generates code
type ProjectGenerator struct {
	Projects map[string]Project
}

// NewProjectGenerator creates a new ProjectGenerator instance
func NewProjectGenerator(projects map[string]Project) *ProjectGenerator {
	return &ProjectGenerator{
		Projects: projects,
	}
}

// ProjectGeneratorData is a template data for project
type ProjectGeneratorData struct {
	Name              string
	ParentPackageName string
}

// Generate generates projects
func (pg *ProjectGenerator) Generate(ctx context.Context) error {
	var errors []error
	for name, project := range pg.Projects {
		projCtx, err := GetContext(ctx)
		if err != nil {
			fmt.Printf("failed to get project %s context: %s\n", name, err.Error())
		}

		err = utils.RemoveFromDirByPatterns(path.Join(name, projCtx.GeneratedFilePrefix+"*"))
		if err != nil {
			return err
		}

		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get project context: %w", err))
			continue
		}
		err = pg.generateProject(projCtx, project, name)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to generate project %s: %w", project.GetType(), err))
		}
	}

	if len(errors) > 0 {
		return utils.NewMultiError("There are errors during project generation", errors...)
	}

	return nil
}

// generateProject generates single project
func (pg *ProjectGenerator) generateProject(projCtx *Context, project Project, name string) error {
	err := project.Init(name, projCtx.ParentPackageName, projCtx.GeneratedFilePrefix)
	if err != nil {
		return fmt.Errorf("failed to init project: %w", err)
	}

	modelgenPlugin := &modelgen.Plugin{
		MutateHook: func(b *modelgen.ModelBuild) *modelgen.ModelBuild {
			return project.MutationHook(b)
		},
		FieldHook: func(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
			return project.ConstraintFieldHook(td, fd, f)
		},
	}

	err = graphql.Generate(
		graphql.PrepareContext(
			context.Background(),
			graphql.GraphqlContext{
				ParentPackageName:   projCtx.ParentPackageName + "/" + name,
				GeneratedFilePrefix: projCtx.GeneratedFilePrefix,
				GeneratorApiOptions: []api.Option{
					api.ReplacePlugin(modelgenPlugin),
				},
				SchemaInjectorHooks: []graphql.InjectorHook{
					func(schema *ast.Schema) generator.SchemaHook {
						err := project.SchemaHook(schema)
						if err != nil {
							panic(err)
						}
						return project.SchemaDocumentHook
					},
				},
				ConfigOverrideCallback: project.ConfigOverride,
			},
		),
		name,
	)

	if err != nil {
		return fmt.Errorf("failed to generate graphql: %w", err)
	}

	return nil
}
