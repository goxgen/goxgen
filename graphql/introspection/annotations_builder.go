package introspection

import (
	"fmt"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

const AnnotationMapTypeName = "XgenAnnotationMap"

func BuildAnnotationIntroHook(schema *ast.Schema, document *ast.SchemaDocument, introValue *map[string]any) error {
	var (
		nameField = &ast.FieldDefinition{
			Name: "name",
			Type: &ast.Type{
				NamedType: "String",
			},
		}
		annotationMapType = &ast.Definition{
			Kind: ast.Object,
			Name: AnnotationMapTypeName,
		}
		annotationField = &ast.FieldDefinition{
			Name: "annotation",
			Type: &ast.Type{
				NamedType: annotationMapType.Name,
			},
		}
	)
	introspectionType := document.Definitions.ForName(TypeName)
	if introspectionType == nil {
		return fmt.Errorf("failed to find XgenIntrospection type")
	}
	introspectionType.Fields = common.AppendFieldIfNotExists(introspectionType.Fields, annotationField)

	query := document.Extensions.ForName("Query")
	if query == nil {
		return fmt.Errorf("failed to find Query type")
	}

	document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, annotationMapType)
	objects := common.GetDefinedObjects(schema)
	annotationValues := make(map[string][]any)
	for _, directive := range schema.Directives {
		if !common.IsXgenDirectiveDefinition(directive) {
			continue
		}
		for _, location := range directive.Locations {
			if location == ast.LocationObject || location == ast.LocationInputObject {
				annotationSingleType := &ast.Definition{
					Kind: ast.Object,
					Name: directive.Name + "AnnotationSingle",
					Fields: []*ast.FieldDefinition{
						nameField,
						{
							Name: "value",
							Type: &ast.Type{
								NamedType: directive.Name,
							},
						},
					},
				}
				document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, annotationSingleType)
				annotationMapType.Fields = common.AppendFieldIfNotExists(annotationMapType.Fields, &ast.FieldDefinition{
					Name: directive.Name,
					Type: ast.NonNullListType(ast.NonNullNamedType(annotationSingleType.Name, nil), nil),
				})
				annotationValues[directive.Name] = make([]any, 0)
			}
		}
	}

	for _, _type := range objects {
		if _type.BuiltIn ||
			_type.Name == "Query" ||
			_type.Name == "Mutation" {
			continue
		}
		for _, directive := range _type.Directives {
			dirValue := make(map[string]any)
			annotationValue := make(map[string]any)
			annotationValue[nameField.Name] = _type.Name
			annotationValue["value"] = &dirValue
			annotationValues[directive.Name] = append(annotationValues[directive.Name], &annotationValue)
			for _, arg := range directive.Arguments {
				val, err := arg.Value.Value(nil)
				if err != nil {
					return fmt.Errorf("failed to get value of %s.%s: %w", _type.Name, directive.Name, err)
				}
				dirValue[arg.Name] = val
			}
		}
	}
	(*introValue)[annotationField.Name] = &annotationValues
	return nil
}
