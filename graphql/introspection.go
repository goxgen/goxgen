package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"strings"
)

// schemaIntrospectionHook is a hook that creates a new schema based on the original schema
func (m *Plugin) schemaIntrospectionHook(schema *ast.Schema) generator.SchemaHook {
	return func(document *ast.SchemaDocument) error {
		var (
			pos       = &ast.Position{Src: &ast.Source{BuiltIn: false}}
			nameField = &ast.FieldDefinition{
				Name: "name",
				Type: &ast.Type{
					NamedType: "String",
				},
			}
			fieldDefType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenFieldDef",
				Position: pos,
				Fields:   []*ast.FieldDefinition{},
			}
			fieldDefsField = &ast.FieldDefinition{
				Name: "definitions",
				Type: &ast.Type{
					NamedType: fieldDefType.Name,
				},
			}
			objectFieldType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenObjectField",
				Position: pos,
				Fields: []*ast.FieldDefinition{
					fieldDefsField,
				},
			}
			objectDefType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenObjectDef",
				Position: pos,
			}
			objDefField = &ast.FieldDefinition{
				Name: "object",
				Type: &ast.Type{
					NamedType: objectDefType.Name,
				},
			}
			objFieldsField = &ast.FieldDefinition{
				Name: "fields",
				Type: ast.NonNullListType(ast.NonNullNamedType(objectFieldType.Name, nil), nil),
			}
			perObjsType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenPerObjects",
				Position: pos,
			}
			perObjectField = &ast.FieldDefinition{
				Name: "_per_object",
				Type: &ast.Type{
					NamedType: perObjsType.Name,
				},
			}
			perDefType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenPerDef",
				Position: pos,
			}
			perDefField = &ast.FieldDefinition{
				Name: "_per_def",
				Type: &ast.Type{
					NamedType: perDefType.Name,
				},
			}
			introspectionType = &ast.Definition{
				Kind:     ast.Object,
				Name:     "XgenIntrospection",
				Position: pos,
				Fields: []*ast.FieldDefinition{
					perObjectField,
					perDefField,
				},
			}
		)

		document.Definitions = generator.AppendDefinitionsIfNotExists(
			document.Definitions,
			introspectionType,
			fieldDefType,
			objectDefType,
			objectFieldType,
			perObjsType,
			perDefType,
		)

		document.Extensions = append(document.Extensions, &ast.Definition{
			Kind:     ast.Object,
			Name:     "Query",
			Position: pos,
			Fields: []*ast.FieldDefinition{
				{
					Name: introspectionQueryField,
					Type: &ast.Type{
						NamedType: introspectionType.Name,
					},
				},
			},
		})

		objects := m.getObjects(schema)
		perDefValues := make(map[string][]any)

		for _, directive := range schema.Directives {
			if !m.isXgenDirectiveDefinition(directive) {
				continue
			}
			newType := m.directiveToType(directive, pos)

			for _, location := range directive.Locations {
				if location == ast.LocationObject || location == ast.LocationInputObject {
					objectDefType.Fields = m.appendFieldIfNotExists(objectDefType.Fields,
						&ast.FieldDefinition{
							Name: directive.Name,
							Type: &ast.Type{
								NamedType: directive.Name,
							},
						},
					)
					perDefSingleType := &ast.Definition{
						Kind:     ast.Object,
						Name:     directive.Name + "PerDefSingle",
						Position: pos,
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
					document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, perDefSingleType)
					perDefType.Fields = m.appendFieldIfNotExists(perDefType.Fields, &ast.FieldDefinition{
						Name: directive.Name,
						Type: ast.NonNullListType(ast.NonNullNamedType(perDefSingleType.Name, nil), nil),
					})
					perDefValues[directive.Name] = make([]any, 0)
				} else if location == ast.LocationFieldDefinition || location == ast.LocationInputFieldDefinition {
					fieldDefType.Fields = m.appendFieldIfNotExists(fieldDefType.Fields,
						&ast.FieldDefinition{
							Name: directive.Name,
							Type: &ast.Type{
								NamedType: directive.Name,
							},
						},
					)
				}
			}

			document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, newType)
		}

		perObjsValues := make(map[string]any)

		for _, _type := range *objects {
			if _type.BuiltIn ||
				_type.Name == "Query" ||
				_type.Name == "Mutation" {
				continue
			}

			objDef := &ast.Definition{
				Kind:     ast.Object,
				Name:     _type.Name + "XgenDef",
				Position: pos,
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
				perDefValue := make(map[string]any)
				perDefValue[nameField.Name] = _type.Name
				perDefValue["value"] = &dirValue
				perDefValues[directive.Name] = append(perDefValues[directive.Name], &perDefValue)
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
			perObjsType.Fields = append(perObjsType.Fields, &ast.FieldDefinition{
				Name: _type.Name,
				Type: &ast.Type{
					NamedType: objDef.Name,
				},
			})
		}

		var introValue = make(map[string]any)
		introValue[perObjectField.Name] = &perObjsValues
		introValue[perDefField.Name] = &perDefValues

		if err := m.saveIntrospectionValuesToFile(introValue); err != nil {
			return fmt.Errorf("failed to save introspection values to file: %w", err)
		}

		return nil
	}
}

func (m *Plugin) directiveToType(directive *ast.DirectiveDefinition, pos *ast.Position) *ast.Definition {
	return &ast.Definition{
		Description: directive.Description,
		Kind:        ast.Object,
		Name:        directive.Name,
		Position:    pos,
		Fields:      m.argsToFields(directive.Arguments),
	}
}

func (m *Plugin) argsToFields(args ast.ArgumentDefinitionList) ast.FieldList {
	fields := ast.FieldList{}
	for _, arg := range args {
		fields = append(fields, &ast.FieldDefinition{
			Name:         arg.Name,
			Description:  arg.Description,
			Type:         arg.Type,
			Position:     arg.Position,
			DefaultValue: arg.DefaultValue,
			Directives:   arg.Directives,
		})
	}
	return fields
}

func (m *Plugin) isXgenDirectiveDefinition(directive *ast.DirectiveDefinition) bool {
	return strings.HasPrefix(directive.Name, "Xgen")
}

func (m *Plugin) saveIntrospectionValuesToFile(values any) error {
	jsonBytes, _ := json.MarshalIndent(values, "", "  ")
	return os.WriteFile(m.introspectionJsonFilePath, jsonBytes, 0644)
}
