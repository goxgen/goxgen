package gqlgen_plugins

import (
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/templates_engine"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
)

type DefTypesInjectorPlugin struct {
	Plugin *Plugin
}

func (m *DefTypesInjectorPlugin) Name() string {
	return "xgen_introspection_types_injector"
}

func (m *DefTypesInjectorPlugin) InjectSourceEarly() *ast.Source {
	schemaRaw, err := os.ReadFile(m.Plugin.introspectionGraphqlFilePath)
	if err != nil {
		panic(fmt.Errorf("unable to open schema: %w", err))
	}
	return &ast.Source{
		BuiltIn: false,
		Name:    m.Plugin.introspectionGraphqlFileName,
		Input:   string(schemaRaw),
	}
}

func (m *DefTypesInjectorPlugin) GenerateCode(data *codegen.Data) error {

	resolverBuild := &ResolverBuild{
		Objects:                   append(data.Objects, data.Inputs...),
		ObjectDirectives:          XgenObjectDirectiveLocation.GetDirectives(data),
		FieldDirectives:           XgenFieldDirectiveLocation.GetDirectives(data),
		PackageName:               m.Plugin.packageName,
		ParentPackageName:         m.Plugin.parentPackageName,
		IntrospectionJsonFileName: m.Plugin.introspectionJsonFileName,
	}

	tbs := &templates_engine.TemplateBundleList{
		&templates_engine.TemplateBundle{
			TemplateDir: "templates/resolver",
			OutputFile:  m.Plugin.resolverFilePath,
			Regenerate:  true,
			FS:          templateFs,
		},
	}
	return tbs.Generate("./", resolverBuild)
}
