package graphql

import (
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/graphql/introspection"
	"github.com/goxgen/goxgen/graphql/resource"
	"github.com/goxgen/goxgen/graphql/validation"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"path"
)

const IntrospectionQueryField = "_xgen_introspection"

type Plugin struct {
	generatedFilePrefix          string
	schemaGeneratorHooks         []InjectorHook
	name                         string
	parentPackageName            string
	introspectionGraphqlFileName string
	introspectionGraphqlFilePath string
	resourcesGraphqlFileName     string
	resourcesGraphqlFilePath     string
	introspectionFileName        string
	introspectionFilePath        string
	commonsFileName              string
	commonsFilePath              string
	introspectionBuildHooks      []introspection.BuilderHook
}

// NewPlugin creates a new plugin
func NewPlugin(
	name string,
	parentPackageName string,
	generatedFilePrefix string,
	schemaGeneratorHooks ...InjectorHook,
) *Plugin {
	p := &Plugin{
		name:                 name,
		parentPackageName:    parentPackageName,
		generatedFilePrefix:  generatedFilePrefix,
		schemaGeneratorHooks: schemaGeneratorHooks,
	}
	p.introspectionBuildHooks = []introspection.BuilderHook{
		introspection.BuildAnnotationIntroHook,
		introspection.BuildPerObjectIntroHook,
		introspection.BuildPerResourceIntroHook,
	}

	p.introspectionGraphqlFileName = p.generatedFilePrefix + "introspection.graphql"
	p.introspectionGraphqlFilePath = path.Join(p.name, consts.GeneratedGqlgenPackageName, p.introspectionGraphqlFileName)

	p.resourcesGraphqlFileName = p.generatedFilePrefix + "resources.graphql"
	p.resourcesGraphqlFilePath = path.Join(p.name, consts.GeneratedGqlgenPackageName, p.resourcesGraphqlFileName)

	p.commonsFileName = p.generatedFilePrefix + "commons.go"
	p.commonsFilePath = path.Join(p.name, p.commonsFileName)

	p.introspectionFileName = p.generatedFilePrefix + "introspection.go"
	p.introspectionFilePath = path.Join(p.name, consts.GeneratedGqlgenPackageName, p.introspectionFileName)

	return p
}

// Name returns the name of the plugin
func (m *Plugin) Name() string {
	return "xgen"
}

func (m *Plugin) Implement(field *codegen.Field) string {
	if field.Name == IntrospectionQueryField {
		return "return " + consts.GeneratedGqlgenPackageName + ".XgenIntrospectionValues()"
	}
	return fmt.Sprintf("panic(fmt.Errorf(\"not implemented: %v - %v\"))", field.GoFieldName, field.Name)
}

type InjectorHook func(schema *ast.Schema) generator.SchemaHook

func (m *Plugin) InjectSourceLate(schema *ast.Schema) *ast.Source {
	var schemaHooks []generator.SchemaHook
	for _, hook := range m.schemaGeneratorHooks {
		schemaHooks = append(schemaHooks, hook(schema))
	}
	introspectionSchemaGenerator := generator.NewSchemaGenerator().
		WithPath(m.introspectionGraphqlFilePath).
		WithSchemaHooks(
			append(
				[]generator.SchemaHook{
					validation.SchemaGeneratorHook(schema, MainDirectiveDefinitionBundle),
					introspection.SchemaGeneratorHook(
						schema,
						IntrospectionQueryField,
						m.introspectionFilePath,
						m.name,
						m.parentPackageName,
						m.introspectionBuildHooks...,
					),
					resource.SchemaGeneratorHook(schema),
				},
				schemaHooks...,
			)...,
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
		Name:    "schema",
		Input:   string(schemaRaw),
	}
}
