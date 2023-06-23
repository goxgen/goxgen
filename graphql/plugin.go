package graphql

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/templates_engine"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"path"
)

//go:embed templates/*
var templateFs embed.FS

const introspectionQueryField = "_xgen_introspection"

type Plugin struct {
	name                         string
	GeneratedFilePrefix          string
	parentPackageName            string
	introspectionGraphqlFileName string
	introspectionGraphqlFilePath string
	resourcesGraphqlFileName     string
	resourcesGraphqlFilePath     string
	introspectionJsonFileName    string
	introspectionJsonFilePath    string
	commonsFileName              string
	commonsFilePath              string
}

type ResolverBuild struct {
	ParentPackageName         string
	PackageName               string
	IntrospectionJsonFileName string
}

// NewPlugin creates a new plugin
func NewPlugin(
	name string,
	parentPackageName string,
	generatedFilePrefix string,
) *Plugin {
	p := &Plugin{
		name:                name,
		parentPackageName:   parentPackageName,
		GeneratedFilePrefix: generatedFilePrefix,
	}

	p.introspectionGraphqlFileName = p.GeneratedFilePrefix + "introspection.graphql"
	p.introspectionGraphqlFilePath = path.Join(p.name, p.introspectionGraphqlFileName)

	p.resourcesGraphqlFileName = p.GeneratedFilePrefix + "resources.graphql"
	p.resourcesGraphqlFilePath = path.Join(p.name, p.resourcesGraphqlFileName)

	p.commonsFileName = p.GeneratedFilePrefix + "commons.go"
	p.commonsFilePath = path.Join(p.name, p.commonsFileName)

	p.introspectionJsonFileName = p.GeneratedFilePrefix + "introspection.json"
	p.introspectionJsonFilePath = path.Join(p.name, p.introspectionJsonFileName)

	return p
}

// Name returns the name of the plugin
func (m *Plugin) Name() string {
	return "xgen"
}

func (m *Plugin) Implement(field *codegen.Field) string {
	if field.Name == introspectionQueryField {
		return "return r.Resolver.XgenIntrospection()"
	}
	return fmt.Sprintf("panic(fmt.Errorf(\"not implemented: %v - %v\"))", field.GoFieldName, field.Name)
}

func (m *Plugin) InjectSourceLate(schema *ast.Schema) *ast.Source {
	introspectionSchemaGenerator := generator.NewSchemaGenerator().
		WithPath(m.introspectionGraphqlFilePath).
		WithSchemaHooks(
			m.schemaIntrospectionHook(schema),
			m.schemaResourcesHook(schema),
		)
	if err := introspectionSchemaGenerator.GenerateOutput(); err != nil {
		panic(err)
	}

	schemaRaw, err := os.ReadFile(m.introspectionGraphqlFilePath)
	if err != nil {
		panic(fmt.Errorf("unable to open schema: %w", err))
	}
	return &ast.Source{
		BuiltIn: false,
		Name:    "xgen-app",
		Input:   string(schemaRaw),
	}
}

func (m *Plugin) GenerateCode(_ *codegen.Data) error {
	resolverBuild := &ResolverBuild{
		PackageName:               m.name,
		ParentPackageName:         m.parentPackageName,
		IntrospectionJsonFileName: m.introspectionJsonFileName,
	}

	tbs := &templates_engine.TemplateBundleList{
		&templates_engine.TemplateBundle{
			TemplateDir: "templates/commons",
			OutputFile:  m.commonsFilePath,
			Regenerate:  true,
			FS:          templateFs,
		},
	}
	return tbs.Generate("./", resolverBuild)
}

func (m *Plugin) getObjects(schema *ast.Schema) *map[string]*ast.Definition {
	objs := make(map[string]*ast.Definition)
	for name, _type := range schema.Types {
		if _type.BuiltIn ||
			_type.Name == "Query" ||
			_type.Name == "Mutation" {
			continue
		}
		objs[name] = _type
	}
	return &objs
}

func (m *Plugin) appendFieldIfNotExists(fields []*ast.FieldDefinition, field *ast.FieldDefinition) []*ast.FieldDefinition {
	for _, f := range fields {
		if f.Name == field.Name {
			return fields
		}
	}
	return append(fields, field)
}
