package goxgen

import (
	"context"
	"github.com/99designs/gqlgen/codegen/config"
	"regexp"
)

// SimpleProject is a default project configuration
type SimpleProject struct {
	name       *string // name of generated code
	outputDir  *string // output directory of generated code
	graphqlURL *string // graphql url

	TemplateBundleList TemplateBundleList
}

// GetName returns project name
func (p *SimpleProject) GetName() *string {
	return p.name
}

// GetOutputDir returns project output directory
func (p *SimpleProject) GetOutputDir() *string {
	return p.outputDir
}

// GetGraphqlURL returns project graphql url
func (p *SimpleProject) GetGraphqlURL() *string {
	return p.graphqlURL
}

func (p *SimpleProject) GqlgenProject() bool {
	return true
}

func (p *SimpleProject) HandleGeneration(ctx context.Context, data *ProjectGeneratorData) error {
	err := p.TemplateBundleList.Generate(PString(p.GetOutputDir()), data)
	if err != nil {
		return err
	}

	err = GenerateProjectGqlgenSet(ctx, p, func(cfg *config.Config) error {
		return nil
	})

	return err
}

// NewProject creates a new SimpleProject instance with default values
func NewProject(name string, options ...ProjectOption) *SimpleProject {

	valid := regexp.MustCompile("^[a-z][a-z0-9_]*$").MatchString(name)
	if !valid {
		panic("project name must be a valid go identifier, \"%s\" provided")
	}

	proj := &SimpleProject{
		name:      StringP(name),
		outputDir: StringP("./"),
		TemplateBundleList: TemplateBundleList{
			{
				TemplateDir: "templates/default-project/handler-templates",
				FS:          templatesFS,
				OutputFile:  "./generated_xgen_project_handlers.go",
				Regenerate:  true,
			},
			{
				TemplateDir: "templates/default-project/graphql-templates",
				FS:          templatesFS,
				OutputFile:  "./generated_xgen_project_graphql.graphql",
				Regenerate:  true,
			},
			{
				TemplateDir: "templates/default-project/resolver",
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
