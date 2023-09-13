package ent

import (
	"embed"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/projects/basic"
	"github.com/goxgen/goxgen/tmpl"
)

// Project is a project configuration for ent projects
type Project struct {
	*basic.Project
}

//go:embed templates/*
var templatesFS embed.FS

type ProjectOption = func(project *Project) error

func WithBasicProjectOption(option basic.ProjectOption) ProjectOption {
	return func(p *Project) error {
		return option(p.Project)
	}
}

func (p *Project) Init(Name string, ParentPackageName string, GeneratedFilePrefix string) error {
	err := p.Project.Init(Name, ParentPackageName, GeneratedFilePrefix)
	if err != nil {
		return err
	}

	tbl := p.StandardTemplateBundleList()

	err = tbl.
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
		Generate(p.Name, p)
	if err != nil {
		return err
	}

	err = projects.RunProjectGoGenCommand(p.Name)
	if err != nil {
		return err
	}

	return (&tmpl.TemplateBundleList{}).Add(
		&tmpl.TemplateBundle{
			TemplateDir: "templates/projects/ent/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  true,
		},
	).Generate(p.Name, p)
}

//
//func (p *Project) PrepareGraphqlGenerationContext(projCtx *projects.Context, data *projects.ProjectGeneratorData) (context.Context, error) {
//	return graphql.PrepareContext(
//		context.Background(),
//		graphql.GraphqlContext{
//			ParentPackageName:   projCtx.ParentPackageName + "/" + data.Name,
//			GeneratedFilePrefix: projCtx.GeneratedFilePrefix,
//			ConfigOverrideCallback: func(cfg *config.Config) error {
//				cfg.AutoBind = append(cfg.AutoBind, projCtx.ParentPackageName+"/"+data.Name+"/ent")
//				cfg.Models.Add("ID", projCtx.ParentPackageName+"/"+data.Name+"/ent/schema/types.UUID")
//				cfg.Models.Add("Node", projCtx.ParentPackageName+"/"+data.Name+"/ent.Noder")
//				cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
//				return nil
//			},
//		},
//	), nil
//}

// New creates a new ent project
func New(option ...ProjectOption) *Project {
	p := &Project{
		Project: basic.NewProject(),
	}
	for _, opt := range option {
		_ = opt(p)
	}
	return p
}
