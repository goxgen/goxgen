package gorm

import (
	"context"
	"embed"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/projects/simple"
	"github.com/goxgen/goxgen/runtime/gorm_initial/generated_gqlgen"
	"github.com/goxgen/goxgen/tmpl"
	"github.com/mitchellh/mapstructure"
	"github.com/vektah/gqlparser/v2/ast"
	"go/types"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

// Project is a default project configuration
type Project struct {
	*simple.Project

	resources []string
}

func (p *Project) MutateHook() modelgen.BuildMutateHook {
	return func(b *modelgen.ModelBuild) *modelgen.ModelBuild {
		for _, model := range b.Models {

			if !slices.Contains(p.resources, model.Name) {
				continue
			}

			for _, f := range model.Fields {
				_, isSlice := f.Type.(*types.Slice)
				if isSlice {
					continue
				}

				pointerType, ok := f.Type.(*types.Pointer)
				if !ok {
					continue
				}

				named, ok := pointerType.Elem().(*types.Named)
				if !ok {
					fmt.Println("Type is not named")
					continue
				}
				fmt.Println("Type is named", named.String(), named.NumMethods())

				typeName := f.Type.String()

				// Check if the type name contains a package prefix and remove it
				if strings.Contains(typeName, ".") {
					components := strings.Split(typeName, ".")
					typeName = components[len(components)-1]
				}

				var typeOfID types.Type = types.Typ[types.Int]

				// Find the ID field of the type
				for _, _m := range b.Models {
					if _m.Name != typeName {
						continue
					}
					for _, _f := range _m.Fields {
						if _f.Name == "id" {
							typeOfID = _f.Type
							break
						}
					}
				}

				if slices.Contains(p.resources, typeName) {
					model.Fields = append(model.Fields, &modelgen.Field{
						Name:   typeName + "ID",
						GoName: typeName + "ID",
						Type:   typeOfID,
					})
				}
			}

			model.Fields = append(model.Fields, &modelgen.Field{
				// gorm.Model
				Type: types.NewPointer(types.NewNamed(
					types.NewTypeName(0, types.NewPackage("gorm.io/gorm", "gorm"), "Model", nil),
					nil,
					nil,
				)),
			})
		}
		return b
	}
}

func (p *Project) getDbConfigFieldDirective(dir *ast.Directive) (*generated_gqlgen.XgenResourceFieldDbConfigInput, error) {
	dbArg := dir.Arguments.ForName("DB")
	dbArgVal, err := dbArg.Value.Value(nil)
	if err != nil {
		return nil, err
	}

	conf := &generated_gqlgen.XgenResourceFieldDbConfigInput{}
	err = mapstructure.Decode(dbArgVal, conf)
	if err != nil {
		fmt.Println("can't decode db config", err)
	}

	return conf, nil
}

func (p *Project) constraintFieldHook() modelgen.FieldMutateHook {
	return func(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
		// Call default hook to proceed standard directives like goField and goTag.
		// You can omit it, if you don't need.
		if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
			return f, err
		}
		gormTag := ``

		c := fd.Directives.ForName(consts.ResourceFieldDirectiveName)
		if c != nil {

			dbConf, err := p.getDbConfigFieldDirective(c)
			if err != nil {
				return nil, err
			}

			if dbConf.Column != nil {
				gormTag += "column:" + *dbConf.Column + ";"
			}

			if dbConf.PrimaryKey != nil && *dbConf.PrimaryKey {
				gormTag += "primaryKey;"
			}

			if dbConf.Unique != nil && *dbConf.Unique {
				gormTag += "unique;"
			}

			if dbConf.Index != nil && *dbConf.Index {
				gormTag += "index;"
			}

			if dbConf.AutoIncrement != nil && *dbConf.AutoIncrement {
				gormTag += "autoIncrement;"
			}

			if dbConf.UniqueIndex != nil && *dbConf.UniqueIndex {
				gormTag += "uniqueIndex;"
			}

			if dbConf.Size != nil {
				gormTag += "size:" + strconv.Itoa(*dbConf.Size) + ";"
			}

			if dbConf.Precision != nil {
				gormTag += "precision:" + strconv.Itoa(*dbConf.Precision) + ";"
			}

			if dbConf.Type != nil {
				gormTag += "type:" + *dbConf.Type + ";"
			}

			if dbConf.Scale != nil {
				gormTag += "scale:" + strconv.Itoa(*dbConf.Scale) + ";"
			}

			if dbConf.NotNull != nil && *dbConf.NotNull {
				gormTag += "not null;"
			}

			if dbConf.AutoIncrementIncrement != nil {
				gormTag += "autoIncrementIncrement:" + strconv.Itoa(*dbConf.AutoIncrementIncrement) + ";"
			}
		}

		f.Tag += fmt.Sprintf(` gorm:"%s"`, gormTag)

		return f, nil
	}
}

type TemplateData struct {
	*simple.TemplateData
	Resources []string
}

func (p *Project) PrepareTemplateData(data *projects.ProjectGeneratorData) *TemplateData {
	return &TemplateData{
		TemplateData: p.Project.PrepareTemplateData(data),
		Resources:    p.resources,
	}
}

func (p *Project) PrepareGraphqlGenerationContext(projCtx *projects.Context, data *projects.ProjectGeneratorData) (context.Context, error) {

	modelgenPlugin := modelgen.Plugin{
		MutateHook: p.MutateHook(),
		FieldHook:  p.constraintFieldHook(),
	}

	gqlCtx, err := p.Project.PrepareGraphqlGenerationContext(projCtx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare graphql generation context: %w", err)
	}

	gqlCtxData, err := graphql.GetGraphqlContext(gqlCtx)
	gqlCtxData.ConfigOverrideCallback = func(cfg *config.Config) error {
		cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
		cfg.Models.Add("ID", "github.com/99designs/gqlgen/graphql.Int")
		return nil
	}
	gqlCtxData.SchemaInjectorHooks = append(gqlCtxData.SchemaInjectorHooks, func(schema *ast.Schema) generator.SchemaHook {
		p.PreserveResources(schema)

		err := p.StandardTemplateBundleList(projCtx).
			Add(
				&tmpl.TemplateBundle{
					TemplateDir: "templates/resolver",
					FS:          templatesFS,
					OutputFile:  "./resolver.go",
					Regenerate:  false,
				},
			).
			Generate(data.Name, p.PrepareTemplateData(data))
		if err != nil {
			panic(err)
		}

		return func(schemaDocument *ast.SchemaDocument) error {
			return nil
		}
	})
	gqlCtxData.GeneratorApiOptions = append(gqlCtxData.GeneratorApiOptions, api.ReplacePlugin(&modelgenPlugin))

	return graphql.PrepareContext(
		context.Background(),
		*gqlCtxData,
	), nil
}

func NewPlugin(option ...projects.ProjectOption) *Project {
	return &Project{
		Project: simple.NewPlugin(option...),
	}
}

type ModelGeneratorData struct {
	PackageName string
	Models      map[string]*ast.Definition
}

func (p *Project) PreserveResources(schema *ast.Schema) {
	allResources := common.GetDefinedObjects(schema, consts.ResourceDirectiveName)
	for _, resource := range allResources {
		p.resources = append(p.resources, resource.Name)
	}
}
