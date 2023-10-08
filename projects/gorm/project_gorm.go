package gorm

import (
	"embed"
	"fmt"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/projects/basic"
	"github.com/goxgen/goxgen/runtime/gorm_initial/generated"
	"github.com/goxgen/goxgen/tmpl"
	"github.com/mitchellh/mapstructure"
	"github.com/vektah/gqlparser/v2/ast"
	"go/types"
	"path"
	"strconv"
	"strings"
)

//go:embed templates/*
var templatesFS embed.FS

// Project is a default project configuration
type Project struct {
	*basic.Project
}

type ProjectOption = func(project *Project) error

func WithBasicProjectOption(option basic.ProjectOption) ProjectOption {
	return func(p *Project) error {
		return option(p.Project)
	}
}

func (p *Project) ModelMutationHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {

	b = p.Project.ModelMutationHook(b)

	for _, model := range b.Models {

		if !p.Resources.TypeExists(model.Name) {
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

			if p.Resources.TypeExists(typeName) {
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

func (p *Project) getDbConfigFieldDirective(dir *ast.Directive) (*generated.XgenResourceFieldDbConfigInput, error) {
	dbArg := dir.Arguments.ForName(consts.SchemaDefFieldDirectiveArgDb)
	dbArgVal, err := dbArg.Value.Value(nil)
	if err != nil {
		return nil, err
	}

	conf := &generated.XgenResourceFieldDbConfigInput{}
	err = mapstructure.Decode(dbArgVal, conf)
	if err != nil {
		fmt.Println("can't decode db config", err)
	}

	return conf, nil
}

func (p *Project) ConstraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	f, err := p.Project.ConstraintFieldHook(td, fd, f)
	if err != nil {
		return f, err
	}

	gormTag := ``

	c := fd.Directives.ForName(consts.SchemaDefDirectiveFieldName)
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

		f.Tag += fmt.Sprintf(` gorm:"%s"`, gormTag)
	}

	return f, nil
}

type TemplateData struct {
	*basic.TemplateData
	Resources map[string]string
}

func (p *Project) PrepareTemplateData() *TemplateData {
	return &TemplateData{
		TemplateData: p.Project.PrepareTemplateData(),
		Resources:    p.Resources,
	}
}

func (p *Project) SchemaHook(schema *ast.Schema) error {
	err := p.Project.SchemaHook(schema)
	if err != nil {
		return err
	}

	return p.StandardTemplateBundleList().
		Add(
			&tmpl.TemplateBundle{
				TemplateDir: "templates/resolver",
				FS:          templatesFS,
				OutputFile:  "./resolver.go",
				Regenerate:  true,
			},
		).
		Add(
			&tmpl.TemplateBundle{
				TemplateDir: "templates/common",
				FS:          templatesFS,
				OutputFile:  path.Join(consts.GeneratedGqlgenPackageName, p.GeneratedFilePrefix+"gorm.go"),
				Regenerate:  true,
			},
		).
		Generate(p.Name, p.PrepareTemplateData())
}

func NewProject(option ...ProjectOption) *Project {
	p := &Project{
		Project: basic.NewProject(),
	}
	for _, opt := range option {
		_ = opt(p)
	}
	return p
}

type ModelGeneratorData struct {
	PackageName string
	Models      map[string]*ast.Definition
}
