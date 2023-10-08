package basic

import (
	"embed"
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/sort"
	"github.com/goxgen/goxgen/tmpl"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/exp/slices"
	"path"
	"strings"
	"text/template"
)

type Resources map[string]string

// Project is a default project configuration
type Project struct {
	typeName                    string    // name of generated code
	testDirectory               string    // name of tests directory
	Resources                   Resources // map of resource name to type name
	ResourceTypeNameToActionMap map[string][]string
	Name                        string
	ParentPackageName           string
	GeneratedFilePrefix         string
}

type TemplateData struct {
	Name              string
	TestsDir          string
	ParentPackageName string

	// Resource related data
	ResourceTypeNameToActionMap map[string][]string
	Resources                   map[string]string

	// Constants
	GeneratedGqlgenPackageName string
}

func (r *Resources) TypeExists(value string) bool {
	for _, v := range *r {
		if v == value {
			return true
		}
	}
	return false
}

func (r *Resources) KeyExists(key string) bool {
	_, exists := (*r)[key]
	return exists
}

func (p *Project) Init(Name string, ParentPackageName string, GeneratedFilePrefix string) error {
	p.Name = Name
	p.ParentPackageName = ParentPackageName
	p.GeneratedFilePrefix = GeneratedFilePrefix
	return nil
}

func (p *Project) GetType() string {
	return p.typeName
}

func (p *Project) TestsDirectory() string {
	return p.testDirectory
}

type ProjectOption = func(project *Project) error

func WithTestDir(dir string) ProjectOption {
	if dir == "" {
		panic("Tests directory cannot be empty")
	}
	return func(p *Project) error {
		p.testDirectory = dir
		return nil
	}
}

//go:embed templates/*
var templatesFS embed.FS

func (p *Project) StandardTemplateBundleList() *tmpl.TemplateBundleList {
	return &tmpl.TemplateBundleList{
		{
			TemplateDir: "templates/handler",
			FS:          templatesFS,
			//OutputFile:  "./" + p.GeneratedFilePrefix + "project_handlers.go",
			OutputFile: path.Join(consts.GeneratedGqlgenPackageName, "server", p.GeneratedFilePrefix+"project_handlers.go"),
			Regenerate: true,
		},
		{
			TemplateDir: "templates/default-tests.yaml",
			FS:          templatesFS,
			//OutputFile:  "./" + p.TestsDirectory() + "/default-tests.yaml",
			OutputFile: path.Join(p.TestsDirectory(), "default-tests.yaml"),
			Regenerate: true,
		},
		{
			TemplateDir: "templates/graphql.config",
			FS:          templatesFS,
			//OutputFile:  "./graphql.config.yml",
			OutputFile: path.Join("graphql.config.yml"),
			Regenerate: true,
		},
		{
			TemplateDir: "templates/resolver",
			FS:          templatesFS,
			//OutputFile:  "./resolver.go",
			OutputFile: path.Join("resolver.go"),
			Regenerate: false,
		},
		{
			TemplateDir: "templates/mapper",
			FS:          templatesFS,
			//OutputFile:  "./" + consts.GeneratedGqlgenPackageName + "/" + p.GeneratedFilePrefix + "mappers.go",
			OutputFile: path.Join(consts.GeneratedGqlgenPackageName, p.GeneratedFilePrefix+"mappers.go"),
			Regenerate: true,
		},
		{
			TemplateDir: "templates/sortable",
			FS:          templatesFS,
			//OutputFile:  "./" + consts.GeneratedGqlgenPackageName + "/" + p.GeneratedFilePrefix + "sortable.go",
			OutputFile: path.Join(consts.GeneratedGqlgenPackageName, p.GeneratedFilePrefix+"sortable.go"),
			Regenerate: true,
			FuncMap: template.FuncMap{
				"singleInputObjectName": sort.ResourceSingleSortInputObjectName,
				"inputObjectName":       sort.ResourceSortInputObjectName,
			},
		},
	}
}

func (p *Project) PrepareTemplateData() *TemplateData {
	return &TemplateData{
		Name:                        p.Name,
		TestsDir:                    p.testDirectory,
		ParentPackageName:           p.ParentPackageName,
		ResourceTypeNameToActionMap: p.ResourceTypeNameToActionMap,
		GeneratedGqlgenPackageName:  consts.GeneratedGqlgenPackageName,
	}
}

func (p *Project) ConstraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	// Call default hook to proceed standard directives like goField and goTag.
	// You can omit it, if you don't need.
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}

	resourceActionDirective := fd.Directives.ForName(consts.SchemaDefDirectiveActionFieldName)
	resourceTagValue := ""
	if resourceActionDirective != nil {
		resourceArg := resourceActionDirective.Arguments.ForName(consts.SchemaDefActionFieldDirectiveArgMapTo)
		if resourceArg != nil {
			resourceFields, err := resourceArg.Value.Value(nil)
			if err != nil {
				return nil, err
			}
			resourceFieldsSlice, ok := resourceFields.([]any)
			if !ok {
				return nil, fmt.Errorf("resource field must be a string")
			}

			resourceFieldsStrs := make([]string, len(resourceFieldsSlice))
			for i, resourceField := range resourceFieldsSlice {
				resourceFieldsStr, ok := resourceField.(string)
				if !ok {
					return nil, fmt.Errorf("resource field must be a string, %T given", resourceField)
				}
				resourceFieldsStrs[i] = resourceFieldsStr
			}

			resourceTagValue += strings.Join(resourceFieldsStrs, ",")
		}

		f.Tag += fmt.Sprintf(` %s:"%s"`, consts.MapToGolangStructTagName, resourceTagValue)
	}

	return f, nil
}

func (p *Project) ModelMutationHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	return b
}

func (p *Project) SchemaHook(schema *ast.Schema) error {
	p.preserveResources(schema)

	return p.StandardTemplateBundleList().
		Generate(
			p.Name,
			p.PrepareTemplateData(),
		)
}

func (p *Project) SchemaDocumentHook(schemaDocument *ast.SchemaDocument) error {
	return nil
}

func (p *Project) ConfigOverride(cfg *config.Config) error {
	err := p.StandardTemplateBundleList().
		Generate(
			p.Name,
			p.PrepareTemplateData(),
		)
	if err != nil {
		return err
	}

	cfg.Models.Add("Map", "github.com/99designs/gqlgen/graphql.Map")
	cfg.Models.Add("ID", "github.com/99designs/gqlgen/graphql.Int")
	return nil
}

// NewProject creates a new Project instance with default values
func NewProject(options ...ProjectOption) *Project {
	proj := &Project{
		testDirectory:               "tests",
		Resources:                   map[string]string{},
		ResourceTypeNameToActionMap: map[string][]string{},
	}
	for _, opt := range options {
		if err := opt(proj); err != nil {
			panic(err)
		}
	}

	return proj
}

func (p *Project) preserveResources(schema *ast.Schema) {
	objects := common.GetDefinedObjects(schema)

	for _, object := range objects {
		resourceDirective := object.Directives.ForName(consts.SchemaDefDirectiveResourceName)
		if resourceDirective != nil {
			resourceName := resourceDirective.Arguments.ForName(consts.SchemaDefResourceDirectiveArgName).Value.Raw
			p.Resources[resourceName] = object.Name
		}

		actionDirectives := append(
			object.Directives.ForNames(consts.SchemaDefDirectiveActionName),
			object.Directives.ForNames(consts.SchemaDefDirectiveListActionName)...,
		)
		for _, actionDirective := range actionDirectives {
			resName := actionDirective.Arguments.ForName(consts.SchemaDefActionDirectiveArgResource)
			if resName == nil {
				panic(fmt.Errorf("resource name is required for %s directive", actionDirective.Name))
			}
			resNameValue, err := resName.Value.Value(nil)
			if err != nil {
				panic(fmt.Errorf("failed to get resource name value: %w", err))
			}
			resNameStr, ok := resNameValue.(string)
			if !ok {
				panic(fmt.Errorf("invalid resource name value: %v", resNameValue))
			}

			resType := common.FindObjectByResourceName(schema, resNameStr)
			if resType == nil {
				panic(fmt.Errorf("mandatory resource %s not found", resNameStr))
			}
			resTypeName := resType.Name
			if !slices.Contains(p.ResourceTypeNameToActionMap[resTypeName], object.Name) {
				p.ResourceTypeNameToActionMap[resTypeName] = append(p.ResourceTypeNameToActionMap[resTypeName], object.Name)
			}
		}
	}
}
