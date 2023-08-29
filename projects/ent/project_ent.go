package ent

import (
	"context"
	"embed"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/projects/simple"
)

// Project is a project configuration for ent projects
type Project struct {
	*simple.Project
}

//go:embed templates/*
var templatesFS embed.FS

func (p *Project) PrepareGraphqlGenerationContext(projCtx *projects.Context, data *projects.ProjectGeneratorData) (context.Context, error) {
	tbl := p.StandardTemplateBundleList(projCtx)

	err := tbl.
		Remove("./resolver.go").
		Add(
			&tmpl.TemplateBundle{
				TemplateDir: "templates/projects/ent/entc",
				FS:          templatesFS,
				OutputFile:  "./ent/entc.go",
				Regenerate:  true,
			},
			&tmpl.TemplateBundle{
				TemplateDir: "templates/projects/ent/gen",
				FS:          templatesFS,
				OutputFile:  "./gen.go",
				Regenerate:  true,
			},
			&tmpl.TemplateBundle{
				TemplateDir: "templates/projects/ent/schema/user",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/user.go",
				Regenerate:  true,
			},
			&tmpl.TemplateBundle{
				TemplateDir: "templates/projects/ent/schema/car",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/car.go",
				Regenerate:  true,
			},
			&tmpl.TemplateBundle{
				TemplateDir: "templates/projects/ent/schema/types",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/types/types.go",
				Regenerate:  true,
			},
		).
		Generate(data.Name, data)
	if err != nil {
		return nil, err
	}

	err = projects.RunProjectGoGenCommand(data.Name)
	if err != nil {
		return nil, err
	}

	err = (&tmpl.TemplateBundleList{}).Add(
		&tmpl.TemplateBundle{
			TemplateDir: "templates/projects/ent/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  true,
		},
	).Generate(data.Name, data)
	if err != nil {
		return nil, err
	}

	return graphql.PrepareContext(
		context.Background(),
		graphql.GraphqlContext{
			ParentPackageName:   projCtx.ParentPackageName + "/" + data.Name,
			GeneratedFilePrefix: projCtx.GeneratedFilePrefix,
			ConfigOverrideCallback: func(cfg *config.Config) error {
				cfg.AutoBind = append(cfg.AutoBind, projCtx.ParentPackageName+"/"+data.Name+"/ent")
				cfg.Models.Add("ID", projCtx.ParentPackageName+"/"+data.Name+"/ent/schema/types.UUID")
				cfg.Models.Add("Node", projCtx.ParentPackageName+"/"+data.Name+"/ent.Noder")
				cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
				return nil
			},
		},
	), nil
}

// NewEntProject creates a new ent project
func NewEntProject(option ...projects.ProjectOption) *Project {
	return &Project{
		Project: simple.New(option...),
	}
}
