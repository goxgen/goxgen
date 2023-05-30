package goxgen

import (
	"context"
	"github.com/99designs/gqlgen/codegen/config"
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
			&TemplateBundle{
				TemplateDir: "templates/ent-templates/entc",
				FS:          templatesFS,
				OutputFile:  "./ent/entc.go",
				Regenerate:  true,
			},
			&TemplateBundle{
				TemplateDir: "templates/ent-templates/gen",
				FS:          templatesFS,
				OutputFile:  "./gen.go",
				Regenerate:  true,
			},
			&TemplateBundle{
				TemplateDir: "templates/ent-templates/schema/user",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/user.go",
				Regenerate:  true,
			},
			&TemplateBundle{
				TemplateDir: "templates/ent-templates/schema/car",
				FS:          templatesFS,
				OutputFile:  "./ent/schema/car.go",
				Regenerate:  true,
			},
			&TemplateBundle{
				TemplateDir: "templates/ent-templates/schema/types",
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

	err = (&TemplateBundleList{}).Add(
		&TemplateBundle{
			TemplateDir: "templates/ent-templates/resolver",
			FS:          templatesFS,
			OutputFile:  "./resolver.go",
			Regenerate:  true,
		},
	).Generate(PString(entp.GetOutputDir()), data)
	if err != nil {
		return err
	}

	err = GenerateProjectGqlgenSet(ctx, entp.SimpleProject, func(cfg *config.Config) error {
		cfg.AutoBind = append(cfg.AutoBind, data.ParentPackageName+"/"+PString(entp.name)+"/ent")
		cfg.Models.Add("ID", data.ParentPackageName+"/"+PString(entp.name)+"/ent/schema/types.UUID")
		cfg.Models.Add("Node", data.ParentPackageName+"/"+PString(entp.name)+"/ent.Noder")
		cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
		return nil
	})

	if err != nil {
		return err
	}

	return err
}

func NewEntProject(name string, option ...ProjectOption) *EntProject {
	return &EntProject{
		SimpleProject: NewProject(name, option...),
	}
}
