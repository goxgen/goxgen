package graphql

import (
	"fmt"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/graphql/inputs"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

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

		for _, object := range *objects {

			resourceActionDirectives := object.Directives.ForNames(directives.XgenResourceActionDirectiveName)
			for _, directive := range resourceActionDirectives {
				xgenDirDef := directives.All.GetInputObjectDirectiveDefinition(directive.Name)
				if xgenDirDef != nil && xgenDirDef.Validate != nil {
					err := xgenDirDef.Validate(directive, object)
					if err != nil {
						return fmt.Errorf("failed to validate object %s: %w", object.Name, err)
					}
				}

				err := m.prepareSchemaField(schema, query, mutation, object, directive)
				if err != nil {
					return fmt.Errorf("failed to prepare field: %w", err)
				}
			}

			resourceListActionDirectives := object.Directives.ForNames(directives.XgenResourceListActionDirectiveName)
			for _, directive := range resourceListActionDirectives {
				xgenDirDef := directives.All.GetInputObjectDirectiveDefinition(directive.Name)
				if xgenDirDef != nil && xgenDirDef.Validate != nil {
					err := xgenDirDef.Validate(directive, object)
					if err != nil {
						return fmt.Errorf("failed to validate object %s: %w", object.Name, err)
					}
				}

				err := m.prepareSchemaField(schema, query, mutation, object, directive)
				if err != nil {
					return fmt.Errorf("failed to prepare field: %w", err)
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

func (m *Plugin) prepareSchemaField(
	schema *ast.Schema,
	query *ast.Definition,
	mutation *ast.Definition,
	object *ast.Definition,
	directive *ast.Directive,
) (err error) {
	resource := directive.Arguments.ForName("Resource").Value.Raw
	isQuery := m.isQueryAction(directive)
	isMutation := m.isMutationAction(directive)

	action := m.getPureActionName(directive)

	returnType, err := m.getObjectSingularTypeByResourceName(schema, directive)
	if err != nil {
		return fmt.Errorf("failed to get object singular type by resource name: %w", err)
	}

	args := ast.ArgumentDefinitionList{
		{
			Name: "input",
			Type: &ast.Type{
				NamedType: object.Name,
			},
		},
	}
	resourceListConfig, _ := directives.GetResourceListConfig(object)

	if resourceListConfig != nil {
		returnType = ast.NonNullListType(returnType, nil)
		if utils.PBool(resourceListConfig.Pagination) {
			args = append(args, &ast.ArgumentDefinition{
				Name: "pagination",
				Type: &ast.Type{
					NamedType: inputs.XgenPaginationInput.Name,
				},
			})
		}
	}

	queryField := &ast.FieldDefinition{
		Name:      resource + "_" + action,
		Arguments: args,
		Type:      returnType,
	}

	if isQuery {
		query.Fields = m.appendFieldIfNotExists(query.Fields, queryField)
	} else if isMutation {
		mutation.Fields = m.appendFieldIfNotExists(mutation.Fields, queryField)
	} else {
		return fmt.Errorf("failed to prepare schema field: unknown action type")
	}
	return nil
}

func (m *Plugin) isQueryAction(directive *ast.Directive) bool {
	resActionEnum := directive.Arguments.ForName("Action").Value.Raw
	return strings.HasSuffix(resActionEnum, "_QUERY")
}

func (m *Plugin) isMutationAction(directive *ast.Directive) bool {
	resActionEnum := directive.Arguments.ForName("Action").Value.Raw
	return strings.HasSuffix(resActionEnum, "_MUTATION")
}

func (m *Plugin) getPureActionName(directive *ast.Directive) string {
	resActionEnum := directive.Arguments.ForName("Action").Value.Raw
	resAction := strings.TrimSuffix(resActionEnum, "_QUERY")
	resAction = strings.TrimSuffix(resAction, "_MUTATION")
	resAction = strings.ToLower(resAction)
	return resAction
}

func (m *Plugin) getXgenResourcePrimaryDirectives(definition *ast.Definition) []*ast.Directive {
	dirs := definition.Directives.ForNames(directives.XgenResourceActionDirectiveName)
	return append(
		dirs,
		definition.Directives.ForNames(directives.XgenResourceListActionDirectiveName)...,
	)
}

func (m *Plugin) getObjectSingularTypeByResourceName(schema *ast.Schema, directive *ast.Directive) (*ast.Type, error) {
	resName := directive.Arguments.ForName("Resource").Value.Raw
	objType := m.findObjectByResourceName(schema, resName)
	if objType == nil {
		return nil, fmt.Errorf("failed to find object for resource %s", resName)
	}

	return &ast.Type{
		NamedType: objType.Name,
	}, nil
}

func (m *Plugin) findObjectByResourceName(schema *ast.Schema, name string) *ast.Definition {
	objects := m.getObjects(schema)
	for _, _type := range *objects {
		directive := _type.Directives.ForName(directives.XgenResourceDirectiveName)

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
