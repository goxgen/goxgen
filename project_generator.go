package goxgen

import (
	"context"
)

// ProjectGenerator is a wrapper for project configuration
// That takes context from xgen and generates code
type ProjectGenerator struct {
	Projects []Project
}

// NewProjectGenerator creates a new ProjectGenerator instance
func NewProjectGenerator(projects ...Project) *ProjectGenerator {
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
	for _, project := range pg.Projects {
		err := pg.generateProject(ctx, project)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateProject generates single project
func (pg *ProjectGenerator) generateProject(ctx context.Context, project Project) error {
	gCtx, err := GetContext(ctx)
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
