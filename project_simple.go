package goxgen

import (
	"context"
	"github.com/goxgen/goxgen/templates_engine"
	"regexp"
)

// SimpleProject is a default project configuration
type SimpleProject struct {
	name               *string // name of generated code
	outputDir          *string // output directory of generated code
	TemplateBundleList templates_engine.TemplateBundleList
}

// GetName returns project name
func (p *SimpleProject) GetName() *string {
	return p.name
}

// GetOutputDir returns project output directory
func (p *SimpleProject) GetOutputDir() *string {
	return p.outputDir
}

// HandleGeneration generates project
func (p *SimpleProject) HandleGeneration(ctx context.Context, data *ProjectGeneratorData) error {
	err := p.TemplateBundleList.Generate(PString(p.GetOutputDir()), data)
	if err != nil {
		return err
	}

	err = GenerateProjectGqlgenSet(ctx, p)

	return err
}

// NewSimpleProject creates a new SimpleProject instance with default values
func NewSimpleProject(name string, options ...ProjectOption) *SimpleProject {

	valid := regexp.MustCompile("^[a-z][a-z0-9_]*$").MatchString(name)
	if !valid {
		panic("project name must be a valid go identifier, \"%s\" provided")
	}

	proj := &SimpleProject{
		name:      StringP(name),
		outputDir: StringP("./"),
		TemplateBundleList: templates_engine.TemplateBundleList{
			{
				TemplateDir: "templates/projects/simple/handler-templates",
				FS:          templatesFS,
				OutputFile:  "./" + GeneratedFilePrefix + "project_handlers.go",
				Regenerate:  true,
			},
			{
				TemplateDir: "templates/projects/simple/resolver",
				FS:          templatesFS,
				OutputFile:  "./resolver.go",
				Regenerate:  false,
			},
		},
	}

	for _, opt := range options {
		if err := opt(proj); err != nil {
			panic(err)
		}
	}

	return proj
}
