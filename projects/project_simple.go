package projects

import (
	"context"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/templates_engine"
)

// SimpleProject is a default project configuration
type SimpleProject struct {
	typeName *string // name of generated code
}

func (p *SimpleProject) GetType() *string {
	return p.typeName
}

func (p *SimpleProject) prepareTemplateBundleList(projCtx *Context) (*templates_engine.TemplateBundleList, error) {
	return &templates_engine.TemplateBundleList{
		{
			TemplateDir: "templates/projects/simple/handler-templates",
			FS:          templatesFS,
			OutputFile:  "./" + projCtx.GeneratedFilePrefix + "project_handlers.go",
			Regenerate:  true,
		},
		{
			TemplateDir: "templates/projects/simple/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  false,
		},
	}, nil
}

func (p *SimpleProject) PrepareGraphqlGenerationContext(projCtx *Context, data *ProjectGeneratorData) (context.Context, error) {

	tbl, err := p.prepareTemplateBundleList(projCtx)
	if err != nil {
		return nil, err
	}

	err = tbl.Generate(data.Name, data)
	if err != nil {
		return nil, err
	}

	return graphql.PrepareContext(context.Background(), graphql.GraphqlContext{
		ParentPackageName:   projCtx.ParentPackageName + "/" + data.Name,
		GeneratedFilePrefix: projCtx.GeneratedFilePrefix,
	}), nil
}

// NewSimpleProject creates a new SimpleProject instance with default values
func NewSimpleProject(options ...ProjectOption) *SimpleProject {
	proj := &SimpleProject{}
	for _, opt := range options {
		if err := opt(proj); err != nil {
			panic(err)
		}
	}

	return proj
}
