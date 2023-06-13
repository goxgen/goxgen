package goxgen

import (
	"context"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/goxgen/goxgen/templates_engine"
)

type EntProject struct {
	*SimpleProject
}

func (entp *EntProject) RunGoGenCommandAfterGenerate() bool {
	return true
}

func (entp *EntProject) HandleGeneration(ctx context.Context, data *ProjectGeneratorData) error {
	err := entp.SimpleProject.TemplateBundleList.
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
		Generate(PString(entp.GetOutputDir()), data)
	if err != nil {
		return err
	}

	err = RunProjectGoGenCommand(entp.SimpleProject)
	if err != nil {
		return err
	}

	err = (&templates_engine.TemplateBundleList{}).Add(
		&templates_engine.TemplateBundle{
			TemplateDir: "templates/projects/ent/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  true,
		},
	).Generate(PString(entp.GetOutputDir()), data)
	if err != nil {
		return err
	}

	err = GenerateProjectGraphqlSet(
		NewGraphqlContext(
			ctx,
			GraphqlContext{
				ConfigOverrideCallback: func(cfg *config.Config) error {
					cfg.AutoBind = append(cfg.AutoBind, data.ParentPackageName+"/"+PString(entp.name)+"/ent")
					cfg.Models.Add("ID", data.ParentPackageName+"/"+PString(entp.name)+"/ent/schema/types.UUID")
					cfg.Models.Add("Node", data.ParentPackageName+"/"+PString(entp.name)+"/ent.Noder")
					cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
					return nil
				},
			},
		),
		entp.SimpleProject,
	)

	if err != nil {
		return err
	}

	return err
}

func NewEntProject(name string, option ...ProjectOption) *EntProject {
	return &EntProject{
		SimpleProject: NewSimpleProject(name, option...),
	}
}
