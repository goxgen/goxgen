package introspection

import (
	"embed"
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/tmpl"
	"github.com/vektah/gqlparser/v2/ast"
)

//go:embed templates/*
var templateFs embed.FS

const TypeName = "XgenIntrospection"

type BuilderHook = func(schema *ast.Schema, document *ast.SchemaDocument, introValue *map[string]any) error

// SchemaGeneratorHook is a hook that creates a new schema based on the original schema
func SchemaGeneratorHook(
	schema *ast.Schema,
	queryFieldName string,
	generatedFilePath string,
	packageName string,
	parentPackageName string,
	introspectionBuildHooks ...BuilderHook) generator.SchemaHook {
	return func(document *ast.SchemaDocument) error {
		var (
			pos               = &ast.Position{Src: &ast.Source{BuiltIn: false}}
			introspectionType = &ast.Definition{
				Kind:     ast.Object,
				Name:     TypeName,
				Position: pos,
				Fields:   []*ast.FieldDefinition{},
			}
		)

		document.Definitions = generator.AppendDefinitionsIfNotExists(
			document.Definitions,
			introspectionType,
		)

		document.Extensions = append(document.Extensions, &ast.Definition{
			Kind:     ast.Object,
			Name:     "Query",
			Position: pos,
			Fields: []*ast.FieldDefinition{
				{
					Name: queryFieldName,
					Type: &ast.Type{
						NamedType: introspectionType.Name,
					},
				},
			},
		})

		for _, directive := range schema.Directives {
			if !common.IsXgenDirectiveDefinition(directive) {
				continue
			}
			newType := common.DirectiveToType(directive, pos)

			document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, newType)
		}

		var introValue = make(map[string]any)

		for _, hook := range introspectionBuildHooks {
			err := hook(schema, document, &introValue)
			if err != nil {
				return fmt.Errorf("failed to run introspection build hook: %w", err)
			}
		}

		if err := saveIntrospectionValuesToFile(
			generatedFilePath,
			introValue,
			packageName,
			parentPackageName,
		); err != nil {
			return fmt.Errorf("failed to save introspection values to file: %w", err)
		}

		return nil
	}
}

type ResolverBuild struct {
	PackageName                string
	GeneratedGqlgenPackageName string
	IntrospectionData          any
}

func saveIntrospectionValuesToFile(generatedFilePath string, values any, packageName string, parentPackageName string) error {
	//jsonBytes, _ := json.MarshalIndent(values, "", "  ")

	data := &ResolverBuild{
		PackageName:                packageName,
		GeneratedGqlgenPackageName: consts.GeneratedGqlgenPackageName,
		IntrospectionData:          values,
	}
	tbs := &tmpl.TemplateBundleList{
		&tmpl.TemplateBundle{
			TemplateDir: "templates",
			OutputFile:  generatedFilePath,
			Regenerate:  true,
			FS:          templateFs,
		},
	}
	return tbs.Generate("./", data)
}
