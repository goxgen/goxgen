package introspection

import (
	"encoding/json"
	"fmt"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
)

const IntrospectionTypeName = "XgenIntrospection"

type BuilderHook = func(schema *ast.Schema, document *ast.SchemaDocument, introValue *map[string]any) error

// SchemaGeneratorHook is a hook that creates a new schema based on the original schema
func SchemaGeneratorHook(schema *ast.Schema, queryFieldName string, generatedJsonFilePath string, introspectionBuildHooks ...BuilderHook) generator.SchemaHook {
	return func(document *ast.SchemaDocument) error {
		var (
			pos = &ast.Position{Src: &ast.Source{BuiltIn: false}}

			introspectionType = &ast.Definition{
				Kind:     ast.Object,
				Name:     IntrospectionTypeName,
				Position: pos,
				Fields:   []*ast.FieldDefinition{
					//perObjectField,
					//perDefField,
				},
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

		if err := saveIntrospectionValuesToFile(generatedJsonFilePath, introValue); err != nil {
			return fmt.Errorf("failed to save introspection values to file: %w", err)
		}

		return nil
	}
}
func saveIntrospectionValuesToFile(generatedJsonFilePath string, values any) error {
	jsonBytes, _ := json.MarshalIndent(values, "", "  ")
	return os.WriteFile(generatedJsonFilePath, jsonBytes, 0644)
}
