package gqlgen_plugins

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/codegen"
	"github.com/goxgen/goxgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"os"
	"path"
	"strings"
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

func appendFieldIfNotExists(fields []*ast.FieldDefinition, field *ast.FieldDefinition) []*ast.FieldDefinition {
	for _, f := range fields {
		if f.Name == field.Name {
			return fields
		}
	}
	return append(fields, field)
}

func (m *Plugin) InjectSourceLate(schema *ast.Schema) *ast.Source {
	schemaGenerator := graphql.SchemaGenerator{
		Path: m.introspectionGraphqlFilePath,
		SchemaHooks: []graphql.SchemaHook{
			func(_schema *ast.Schema) error {
				pos := &ast.Position{Src: &ast.Source{BuiltIn: false}}

				introspectionType := &ast.Definition{
					Kind:     ast.Object,
					Name:     "XgenIntrospection",
					Position: pos,
				}
				fieldDefType := &ast.Definition{
					Kind:     ast.Object,
					Name:     "XgenFieldDef",
					Position: pos,
				}
				objectDefType := &ast.Definition{
					Kind:     ast.Object,
					Name:     "XgenObjectDef",
					Position: pos,
				}

				_schema.AddTypes(fieldDefType, objectDefType, introspectionType)

				for _, directive := range schema.Directives {
					if !strings.HasPrefix(directive.Name, "Xgen") {
						continue
					}
					newType := &ast.Definition{
						Description: directive.Description,
						Kind:        ast.Object,
						Name:        directive.Name,
						Position:    pos,
					}

					for _, location := range directive.Locations {
						if location == ast.LocationObject || location == ast.LocationInputObject {
							objectDefType.Fields = appendFieldIfNotExists(objectDefType.Fields,
								&ast.FieldDefinition{
									Name: directive.Name,
									Type: &ast.Type{
										NamedType: directive.Name,
									},
								},
							)
						} else if location == ast.LocationFieldDefinition || location == ast.LocationInputFieldDefinition {
							fieldDefType.Fields = appendFieldIfNotExists(fieldDefType.Fields,
								&ast.FieldDefinition{
									Name: directive.Name,
									Type: &ast.Type{
										NamedType: directive.Name,
									},
								},
							)
						}
					}

					for _, arg := range directive.Arguments {
						newType.Fields = append(newType.Fields, &ast.FieldDefinition{
							Name:         arg.Name,
							Description:  arg.Description,
							Type:         arg.Type,
							Position:     arg.Position,
							DefaultValue: arg.DefaultValue,
							Directives:   arg.Directives,
						})
					}

					_schema.AddTypes(newType)
				}

				values := make(map[string]any)
				objDefField := &ast.FieldDefinition{
					Name: "_object_def",
					Type: &ast.Type{
						NamedType: objectDefType.Name,
					},
				}

				for _, _type := range schema.Types {
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
						},
					}

					objDefValue := make(map[string]any)
					objDefFieldValue := make(map[string]any)
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

					for _, field := range _type.Fields {
						fieldDef := &ast.FieldDefinition{
							Name: field.Name,
							Type: &ast.Type{
								NamedType: fieldDefType.Name,
							},
						}
						objDef.Fields = append(objDef.Fields, fieldDef)

						fieldDefValue := make(map[string]any)
						for _, directive := range field.Directives {
							dirValue := make(map[string]any)
							for _, arg := range directive.Arguments {
								val, err := arg.Value.Value(nil)
								if err != nil {
									return fmt.Errorf("failed to get value of %s.%s.%s: %w", _type.Name, field.Name, arg.Name, err)
								}
								dirValue[arg.Name] = val
							}
							fieldDefValue[directive.Name] = &dirValue
						}
						objDefValue[fieldDef.Name] = &fieldDefValue
					}

					values[_type.Name] = &objDefValue

					_schema.AddTypes(objDef)
					introspectionType.Fields = append(introspectionType.Fields, &ast.FieldDefinition{
						Name: _type.Name,
						Type: &ast.Type{
							NamedType: objDef.Name,
						},
					})
				}

				jsonBytes, _ := json.MarshalIndent(values, "", "  ")
				_ = os.WriteFile(m.introspectionJsonFilePath, jsonBytes, 0644)

				return nil
			},
		},
		Footer: `extend type Query { _xgen_introspection: XgenIntrospection }`,
	}

	if err := schemaGenerator.GenerateOutput(); err != nil {
		panic(err)
	}

	return nil
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
