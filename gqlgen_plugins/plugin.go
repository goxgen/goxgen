package gqlgen_plugins

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/templates_engine"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFs embed.FS

type Plugin struct {
	Directory           string
	GeneratedFilePrefix string
	packageName         string
	parentPackageName   string

	introspectionGraphqlFileName string
	introspectionGraphqlFilePath string

	introspectionJsonFileName string
	introspectionJsonFilePath string

	directivesFileName string
	directivesFilePath string

	resolverFileName string
	resolverFilePath string
}

func NewPlugin(
	directory string,
	packageName string,
	parentPackageName string,
	generatedFilePrefix string,
) (*Plugin, *DefTypesInjectorPlugin) {
	p := &Plugin{
		Directory:           directory,
		packageName:         packageName,
		parentPackageName:   parentPackageName,
		GeneratedFilePrefix: generatedFilePrefix,
	}

	p.introspectionGraphqlFileName = p.GeneratedFilePrefix + "introspection.graphql"
	p.introspectionGraphqlFilePath = path.Join(p.Directory, p.introspectionGraphqlFileName)

	p.resolverFileName = p.GeneratedFilePrefix + "introspection.graphql.resolver.go"
	p.resolverFilePath = path.Join(p.Directory, p.resolverFileName)

	p.introspectionJsonFileName = p.GeneratedFilePrefix + "introspection.json"
	p.introspectionJsonFilePath = path.Join(p.Directory, p.introspectionJsonFileName)

	p.directivesFileName = p.GeneratedFilePrefix + "directives.graphql"
	p.directivesFilePath = path.Join(p.Directory, p.directivesFileName)

	return p, &DefTypesInjectorPlugin{Plugin: p}
}

func (m *Plugin) Name() string {
	return "xgen"
}

type XgenDirectiveLocation struct {
	Name      string
	Locations []ast.DirectiveLocation
}

type XgenDirectiveLocations []*XgenDirectiveLocation

var XgenObjectDirectiveLocation = XgenDirectiveLocation{
	Name: "Object",
	Locations: []ast.DirectiveLocation{
		ast.LocationObject,
		ast.LocationInputObject,
	},
}

var XgenFieldDirectiveLocation = XgenDirectiveLocation{
	Name: "Field",
	Locations: []ast.DirectiveLocation{
		ast.LocationFieldDefinition,
		ast.LocationInputFieldDefinition,
	},
}

func (m *Plugin) InjectSourceEarly() *ast.Source {
	tb := &templates_engine.TemplateBundle{
		TemplateDir: "templates/directives",
		OutputFile:  m.directivesFilePath,
		Regenerate:  true,
		FS:          templateFs,
	}
	err := tb.Generate("./", nil)
	if err != nil {
		panic(err)
	}
	var schemaRaw []byte
	schemaRaw, err = os.ReadFile(m.directivesFilePath)
	if err != nil {
		panic(fmt.Errorf("unable to open schema: %w", err))
	}

	return &ast.Source{
		BuiltIn: false,
		Name:    m.directivesFileName,
		Input:   string(schemaRaw),
	}
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	resolverBuild := &ResolverBuild{
		Objects:                   append(data.Objects, data.Inputs...),
		ObjectDirectives:          XgenObjectDirectiveLocation.GetDirectives(data),
		FieldDirectives:           XgenFieldDirectiveLocation.GetDirectives(data),
		PackageName:               m.packageName,
		ParentPackageName:         m.parentPackageName,
		IntrospectionJsonFileName: m.introspectionJsonFileName,
	}

	fm := template.FuncMap{
		"fieldName": func(field *codegen.Field) string {
			return strings.TrimPrefix(field.Name, "Xgen")
		},
		"directiveName": func(dir *codegen.Directive) string {
			return strings.TrimPrefix(dir.Name, "Xgen")
		},
	}

	tbs := &templates_engine.TemplateBundleList{
		&templates_engine.TemplateBundle{
			TemplateDir: "templates/types-definitions",
			OutputFile:  m.introspectionGraphqlFilePath,
			Regenerate:  true,
			FS:          templateFs,
			FuncMap:     fm,
		},
		&templates_engine.TemplateBundle{
			TemplateDir: "templates/definitions",
			OutputFile:  m.introspectionJsonFilePath,
			Regenerate:  true,
			FS:          templateFs,
			FuncMap:     fm,
		},
	}
	return tbs.Generate("./", resolverBuild)
}

func isXgenDirective(dir *codegen.Directive) bool {
	return strings.HasPrefix(dir.Name, "Xgen")
}

func (xl XgenDirectiveLocation) GetDirectives(data *codegen.Data) *codegen.DirectiveList {
	directives := codegen.DirectiveList{}

	for _, location := range xl.Locations {
		for _, dir := range data.AllDirectives.LocationDirectives(string(location)) {
			if isXgenDirective(dir) {
				directives[dir.Name] = dir
			}
		}
	}

	return &directives
}

func (dls *XgenDirectiveLocations) GetAllDirectives(data *codegen.Data) *codegen.DirectiveList {
	directives := codegen.DirectiveList{}

	for _, dl := range *dls {
		for _, dir := range *dl.GetDirectives(data) {
			directives[dir.Name] = dir
		}
	}

	return &directives
}

type ResolverBuild struct {
	Objects []*codegen.Object
	Inputs  []*codegen.Object

	ObjectDirectives          *codegen.DirectiveList
	FieldDirectives           *codegen.DirectiveList
	ParentPackageName         string
	PackageName               string
	IntrospectionJsonFileName string
}
