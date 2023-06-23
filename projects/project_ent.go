package projects

import (
	"context"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/templates_engine"
)

// EntProject is a project configuration for ent projects
type EntProject struct {
	*SimpleProject
}

func (entp *EntProject) PrepareGraphqlGenerationContext(projCtx *Context, data *ProjectGeneratorData) (context.Context, error) {
	tbl, err := entp.SimpleProject.prepareTemplateBundleList(projCtx)
	if err != nil {
		return nil, err
	}

	err = tbl.
		Remove("./resolver.go").
		Add(
			&templates_engine.TemplateBundle{
				TemplateDir: "templates/projects/ent/entc",
				FS:          templatesFS,
				OutputFile:  "./ent/entc.go",
				Regenerate:  true,
			},
			&templates_engine.TemplateBundle{
				TemplateDir: "templates/projects/ent/gen",
				FS:          templatesFS,
				OutputFile:  "./gen.go",
				Regenerate:  true,
			},
			&templates_engine.TemplateBundle{
				TemplateDir: "templates/projects/ent/schema/user",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/user.go",
				Regenerate:  true,
			},
			&templates_engine.TemplateBundle{
				TemplateDir: "templates/projects/ent/schema/car",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/car.go",
				Regenerate:  true,
			},
			&templates_engine.TemplateBundle{
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

	err = RunProjectGoGenCommand(data.Name)
	if err != nil {
		return nil, err
	}

	err = (&templates_engine.TemplateBundleList{}).Add(
		&templates_engine.TemplateBundle{
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
func NewEntProject(option ...ProjectOption) *EntProject {
	return &EntProject{
		SimpleProject: NewSimpleProject(option...),
	}
}
