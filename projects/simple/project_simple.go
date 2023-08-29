package simple

import (
	"context"
	"embed"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/tmpl"
)

// Project is a default project configuration
type Project struct {
	typeName *string // name of generated code
}

type TemplateData struct {
	*projects.ProjectGeneratorData
	GeneratedGqlgenPackageName string
}

func (p *Project) GetType() *string {
	return p.typeName
}

//go:embed templates/*
var templatesFS embed.FS

func (p *Project) StandardTemplateBundleList(projCtx *projects.Context) *tmpl.TemplateBundleList {
	return &tmpl.TemplateBundleList{
		{
			TemplateDir: "templates/handler-templates",
			FS:          templatesFS,
			OutputFile:  "./" + projCtx.GeneratedFilePrefix + "project_handlers.go",
			Regenerate:  true,
		},
		{
			TemplateDir: "templates/graphql.config",
			FS:          templatesFS,
			OutputFile:  "./graphql.config.yml",
			Regenerate:  true,
		},
		{
			TemplateDir: "templates/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  false,
		},
	}
}

func (p *Project) PrepareTemplateData(data *projects.ProjectGeneratorData) *TemplateData {
	return &TemplateData{
		ProjectGeneratorData:       data,
		GeneratedGqlgenPackageName: consts.GeneratedGqlgenPackageName,
	}
}

func (p *Project) PrepareGraphqlGenerationContext(projCtx *projects.Context, data *projects.ProjectGeneratorData) (context.Context, error) {

	err := p.StandardTemplateBundleList(projCtx).Generate(
		data.Name,
		p.PrepareTemplateData(data),
	)
	if err != nil {
		return nil, err
	}

	return graphql.PrepareContext(
		context.Background(),
		graphql.GraphqlContext{
			ParentPackageName:   projCtx.ParentPackageName + "/" + data.Name,
			GeneratedFilePrefix: projCtx.GeneratedFilePrefix,
			ConfigOverrideCallback: func(cfg *config.Config) error {
				cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
				cfg.Models.Add("ID", "github.com/99designs/gqlgen/graphql.Int")
				return nil
			},
		},
	), nil
}

// New creates a new Project instance with default values
func New(options ...projects.ProjectOption) *Project {
	proj := &Project{}
	for _, opt := range options {
		if err := opt(proj); err != nil {
			panic(err)
		}
	}

	return proj
}
