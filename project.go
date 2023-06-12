package goxgen

import (
	"context"
	"github.com/goxgen/goxgen/utils"
)

// Project is a project configuration
type Project interface {
	GetName() *string
	GetOutputDir() *string
	HandleGeneration(ctx context.Context, data *ProjectGeneratorData) error
}

type ProjectWithCustomTemplateData interface {
	Project
	PrepareCustomTemplateData(ctx context.Context, data *ProjectGeneratorData) error
}

func RunProjectGoGenCommand(project Project) error {
	outputDir := PString(project.GetOutputDir())
	return utils.ExecCommand(outputDir, "go", "generate")
}
