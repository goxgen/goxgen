package introspection

import (
	"fmt"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

const ObjectDefinitionName = "XgenObjectDefinition"

func BuildPerObjectIntroHook(schema *ast.Schema, document *ast.SchemaDocument, introValue *map[string]any) error {
	var (
		nameField = &ast.FieldDefinition{
			Name: "name",
			Type: &ast.Type{
				NamedType: "String",
			},
		}
		objectMapType = &ast.Definition{
			Kind: ast.Object,
			Name: "XgenObjectMap",
		}
		fieldDefType = &ast.Definition{
			Kind:   ast.Object,
			Name:   "XgenFieldDef",
			Fields: []*ast.FieldDefinition{},
		}
		fieldDefsField = &ast.FieldDefinition{
			Name: "definition",
			Type: &ast.Type{
				NamedType: fieldDefType.Name,
			},
		}
		objectField = &ast.FieldDefinition{
			Name: "object",
			Type: &ast.Type{
				NamedType: objectMapType.Name,
			},
		}
		objectFieldType = &ast.Definition{
			Kind: ast.Object,
			Name: "XgenObjectField",
			Fields: []*ast.FieldDefinition{
				nameField,
				fieldDefsField,
			},
		}
		objectDefType = &ast.Definition{
			Kind: ast.Object,
			Name: ObjectDefinitionName,
		}
		objDefField = &ast.FieldDefinition{
			Name: "object",
			Type: &ast.Type{
				NamedType: objectDefType.Name,
			},
		}
		objFieldsField = &ast.FieldDefinition{
			Name: "field",
			Type: ast.NonNullListType(ast.NonNullNamedType(objectFieldType.Name, nil), nil),
		}
	)

	document.Definitions = generator.AppendDefinitionsIfNotExists(
		document.Definitions,
		fieldDefType,
		objectDefType,
		objectFieldType,
		objectMapType,
	)

	introspectionType := document.Definitions.ForName(TypeName)
	if introspectionType == nil {
		return fmt.Errorf("failed to find XgenIntrospection type")
	}
	introspectionType.Fields = common.AppendFieldIfNotExists(introspectionType.Fields, objectField)

	query := document.Extensions.ForName("Query")
	if query == nil {
		return fmt.Errorf("failed to find Query type")
	}

	for _, directive := range schema.Directives {
		if !common.IsXgenDirectiveDefinition(directive) {
			continue
		}

		for _, location := range directive.Locations {
			if location == ast.LocationObject || location == ast.LocationInputObject {
				objectDefType.Fields = common.AppendFieldIfNotExists(objectDefType.Fields,
					&ast.FieldDefinition{
						Name: directive.Name,
						Type: &ast.Type{
							NamedType: directive.Name,
						},
					},
				)

			} else if location == ast.LocationFieldDefinition || location == ast.LocationInputFieldDefinition {
				fieldDefType.Fields = common.AppendFieldIfNotExists(fieldDefType.Fields,
					&ast.FieldDefinition{
						Name: directive.Name,
						Type: &ast.Type{
							NamedType: directive.Name,
						},
					},
				)
			}
		}

	}

	perObjsValues := make(map[string]any)

	objects := common.GetDefinedObjects(schema)

	for _, _type := range objects {
		if _type.BuiltIn ||
			_type.Name == "Query" ||
			_type.Name == "Mutation" {
			continue
		}

		objDef := &ast.Definition{
			Kind: ast.Object,
			Name: _type.Name + "XgenDef",
			Fields: []*ast.FieldDefinition{
				objDefField,
				objFieldsField,
			},
		}

		objDefValue := make(map[string]any)
		objDefFieldValue := make(map[string]any)
		objFieldsFieldValue := make([]map[string]any, 0)

		for _, directive := range _type.Directives {
			dirValue := make(map[string]any)
			for _, arg := range directive.Arguments {
				val, err := arg.Value.Value(nil)
				if err != nil {
					return fmt.Errorf("failed to get value of %s.%s: %w", _type.Name, directive.Name, err)
				}
				dirValue[arg.Name] = val
			}
			objDefFieldValue[directive.Name] = &dirValue
		}
		objDefValue[objDefField.Name] = &objDefFieldValue
		objDefValue[objFieldsField.Name] = &objFieldsFieldValue

		for _, field := range _type.Fields {
			fieldDef := make(map[string]any)
			fieldDef[nameField.Name] = field.Name
			fieldDefs := make(map[string]any)
			fieldDef[fieldDefsField.Name] = &fieldDefs
			for _, directive := range field.Directives {
				dirValue := make(map[string]any)
				for _, arg := range directive.Arguments {
					val, err := arg.Value.Value(nil)
					if err != nil {
						return fmt.Errorf("failed to get value of %s.%s.%s: %w", _type.Name, field.Name, arg.Name, err)
					}
					dirValue[arg.Name] = val
				}
				fieldDefs[directive.Name] = &dirValue
			}
			objFieldsFieldValue = append(objFieldsFieldValue, fieldDef)
		}

		perObjsValues[_type.Name] = &objDefValue

		document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, objDef)
		objectMapType.Fields = append(objectMapType.Fields, &ast.FieldDefinition{
			Name: _type.Name,
			Type: &ast.Type{
				NamedType: objDef.Name,
			},
		})
	}

	(*introValue)[objectField.Name] = &perObjsValues

	return nil

}
