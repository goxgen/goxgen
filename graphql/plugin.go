package graphql

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/templates_engine"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"path"
	"strings"
)

//go:embed templates/*
var templateFs embed.FS

const introspectionQueryField = "_xgen_introspection"

type Plugin struct {
	name                string
	GeneratedFilePrefix string
	parentPackageName   string

	introspectionGraphqlFileName string
	introspectionGraphqlFilePath string

	resourcesGraphqlFileName string
	resourcesGraphqlFilePath string

	introspectionJsonFileName string
	introspectionJsonFilePath string

	commonsFileName string
	commonsFilePath string
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

func (m *Plugin) isXgenDirectiveDefinition(directive *ast.DirectiveDefinition) bool {
	return strings.HasPrefix(directive.Name, "Xgen")
}

func (m *Plugin) isXgenDirective(directive *ast.Directive) bool {
	return strings.HasPrefix(directive.Name, "Xgen")
}

// schemaIntrospectionHook is a hook that creates a new schema based on the original schema
func (m *Plugin) schemaIntrospectionHook(schema *ast.Schema) generator.SchemaHook {
	return func(_document *ast.SchemaDocument) error {
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

		_document.Definitions = generator.AppendDefinitionsIfNotExists(
			_document.Definitions,
			introspectionType,
			fieldDefType,
			objectDefType,
			objectFieldType,
			perObjsType,
			perDefType,
		)

		_document.Extensions = append(_document.Extensions, &ast.Definition{
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
					_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, perDefSingleType)
					perDefType.Fields = m.appendFieldIfNotExists(perDefType.Fields, &ast.FieldDefinition{
						Name: directive.Name,
						Type: &ast.Type{
							NamedType: "[" + perDefSingleType.Name + "]",
						},
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

			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, newType)
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

			_document.Definitions = generator.AppendDefinitionsIfNotExists(_document.Definitions, objDef)
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

func (m *Plugin) findObjectByResourceName(schema *ast.Schema, name string) *ast.Definition {
	objects := m.getObjects(schema)
	for _, _type := range *objects {
		directive := _type.Directives.ForName(directives.XgenResource.Name)

		if directive == nil {
			continue
		}

		resNameArg := directive.Arguments.ForName("Name")
		if resNameArg != nil && resNameArg.Value.Raw == name {
			return _type
		}

	}
	return nil
}

func (m *Plugin) schemaResourcesHook(schema *ast.Schema) generator.SchemaHook {
	return func(document *ast.SchemaDocument) error {
		objects := m.getObjects(schema)
		query := &ast.Definition{
			Kind:   ast.Object,
			Name:   "Query",
			Fields: []*ast.FieldDefinition{},
		}
		mutation := &ast.Definition{
			Kind:   ast.Object,
			Name:   "Mutation",
			Fields: []*ast.FieldDefinition{},
		}

		for _, _type := range *objects {
			resQueryDirectives := _type.Directives.ForNames(directives.XgenResourceAction.Name)
			resQueryDirectives = append(resQueryDirectives, _type.Directives.ForNames(directives.XgenResourceListAction.Name)...)

			resourceListConfig, _ := directives.GetResourceListConfig(_type)

			for _, directive := range resQueryDirectives {
				resName := directive.Arguments.ForName("Resource").Value.Raw
				resActionEnum := directive.Arguments.ForName("Action").Value.Raw

				isQuery := strings.HasSuffix(resActionEnum, "_QUERY")
				isMutation := strings.HasSuffix(resActionEnum, "_MUTATION")

				resAction := strings.TrimSuffix(resActionEnum, "_QUERY")
				resAction = strings.TrimSuffix(resAction, "_MUTATION")
				resAction = strings.ToLower(resAction)

				objType := m.findObjectByResourceName(schema, resName)
				if objType == nil {
					return fmt.Errorf("failed to find object for resource %s", resName)
				}

				args := ast.ArgumentDefinitionList{
					{
						Name: "input",
						Type: &ast.Type{
							NamedType: _type.Name,
						},
					},
				}
				returnType := &ast.Type{
					NamedType: objType.Name,
				}
				if resourceListConfig != nil {
					returnType = ast.NonNullListType(returnType, nil)
					if utils.PBool(resourceListConfig.Pagination) {
						//args = append(args, &ast.ArgumentDefinition{
						//	Name: "pagination",
						//	Type: &ast.Type{
						//		NamedType: inputs.XgenPaginationInput.Name,
						//	},
						//})
					}
				}
				queryField := &ast.FieldDefinition{
					Name:      resName + "_" + resAction,
					Arguments: args,
					Type:      returnType,
				}
				if isQuery {
					query.Fields = m.appendFieldIfNotExists(query.Fields, queryField)
				} else if isMutation {
					mutation.Fields = m.appendFieldIfNotExists(mutation.Fields, queryField)
				} else {
					return fmt.Errorf("unknown action type %s", resAction)
				}

			}
		}

		if len(query.Fields) > 0 {
			document.Extensions = append(document.Extensions, query)
		}

		if len(mutation.Fields) > 0 {
			document.Extensions = append(document.Extensions, mutation)
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

func (m *Plugin) saveIntrospectionValuesToFile(values any) error {
	jsonBytes, _ := json.MarshalIndent(values, "", "  ")
	return os.WriteFile(m.introspectionJsonFilePath, jsonBytes, 0644)
}

func (m *Plugin) appendFieldIfNotExists(fields []*ast.FieldDefinition, field *ast.FieldDefinition) []*ast.FieldDefinition {
	for _, f := range fields {
		if f.Name == field.Name {
			return fields
		}
	}
	return append(fields, field)
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
