package projects

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/utils"
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
			errors = append(errors, fmt.Errorf("failed to generate project %s: %w", utils.PString(project.GetType()), err))
		}
	}

	if len(errors) > 0 {
		return utils.NewMultiError("There are errors during project generation", errors...)
	}

	return nil
}

// generateProject generates single project
func (pg *ProjectGenerator) generateProject(projCtx *Context, project Project, name string) error {

	genData := &ProjectGeneratorData{
		ParentPackageName: projCtx.ParentPackageName,
		Name:              name,
	}

	if projectWithTemplateData, ok := project.(ProjectWithCustomTemplateData); ok {
		err := projectWithTemplateData.PrepareCustomTemplateData(projCtx, genData)
		if err != nil {
			return err
		}
	}

	graphqlCtx, err := project.PrepareGraphqlGenerationContext(projCtx, genData)
	if err != nil {
		return fmt.Errorf("failed to prepare graphql generation: %w", err)
	}
	err = graphql.Generate(
		graphqlCtx,
		name,
	)

	if err != nil {
		return fmt.Errorf("failed to generate graphql: %w", err)
	}

	return nil
}
