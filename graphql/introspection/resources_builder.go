package introspection

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

func BuildPerResourceIntroHook(schema *ast.Schema, document *ast.SchemaDocument, introValue *map[string]any) error {
	resourceDirective, exists := schema.Directives[consts.ResourceDirectiveName]
	if !exists {
		return fmt.Errorf("failed to find %s directive", consts.ResourceDirectiveName)
	}
	resourceActionDirective, exists := schema.Directives[consts.ResourceActionDirectiveName]
	if !exists {
		return fmt.Errorf("failed to find %s directive", consts.ResourceDirectiveName)
	}
	var (
		resourceMapType = &ast.Definition{
			Kind: ast.Object,
			Name: "XgenResourceMap",
		}
		resourceQueryField = &ast.FieldDefinition{
			Name: "resource",
			Type: &ast.Type{
				NamedType: resourceMapType.Name,
			},
		}
		resourcePropertyType = &ast.Definition{
			Kind:   ast.Object,
			Name:   "XgenResourceProperty",
			Fields: common.ArgsToFields(resourceDirective.Arguments),
		}
		actionType = &ast.Definition{
			Kind:   ast.Object,
			Name:   "XgenResourceAction",
			Fields: common.ArgsToFields(resourceActionDirective.Arguments),
		}
		resourceDefinitionType = &ast.Definition{
			Kind: ast.Object,
			Name: "XgenResourceDefinition",
			Fields: []*ast.FieldDefinition{
				{
					Name: "objectName",
					Type: &ast.Type{
						NamedType: "String",
					},
				},
				{
					Name: "properties",
					Type: &ast.Type{
						NamedType: resourcePropertyType.Name,
					},
				},
				{
					Name: "actions",
					Type: ast.NonNullListType(&ast.Type{
						NamedType: actionType.Name,
					}, nil),
				},
			},
		}
	)

	objects := common.GetDefinedObjects(schema)

	resourceValue := make(map[string]any)

	for _, object := range objects {
		resourceDirective := object.Directives.ForName(consts.ResourceDirectiveName)
		if resourceDirective == nil {
			continue
		}

		resourceNameArg := resourceDirective.Arguments.ForName("Name")
		if resourceNameArg == nil {
			return fmt.Errorf("failed to find Name argument in %s directive", resourceDirective.Name)
		}
		resourceName, err := resourceNameArg.Value.Value(nil)
		if err != nil {
			return fmt.Errorf("failed to get value of Name argument in %s directive", resourceDirective.Name)
		}
		resourceNameStr, ok := resourceName.(string)
		if !ok {
			return fmt.Errorf("failed to cast Name argument value to string in %s directive", resourceDirective.Name)
		}

		resourceMapType.Fields = common.AppendFieldIfNotExists(
			resourceMapType.Fields,
			&ast.FieldDefinition{
				Name: resourceNameStr,
				Type: &ast.Type{
					NamedType: resourceDefinitionType.Name,
				},
			},
		)
		singleResourceValue := make(map[string]any)
		singleResourceValue["objectName"] = object.Name

		props := make(map[string]any)
		singleResourceValue["properties"] = &props
		actions := make([]map[string]any, 0)
		singleResourceValue["actions"] = &actions

		for _, arg := range resourceDirective.Arguments {
			val, err := arg.Value.Value(nil)
			if err != nil {
				return fmt.Errorf("failed to get value of %s argument in %s directive", arg.Name, resourceDirective.Name)
			}
			props[arg.Name] = val
		}

		for _, _object := range objects {
			actionDirs := directives.GetResourceActionDirectives(_object)
			for _, actionDirective := range actionDirs {
				_resourceNameArg := actionDirective.Arguments.ForName("Resource")
				if _resourceNameArg == nil {
					return fmt.Errorf("failed to find Resource argument in %s directive", actionDirective.Name)
				}
				_resourceName, err := _resourceNameArg.Value.Value(nil)
				if err != nil {
					return fmt.Errorf("failed to get value of Resource argument in %s directive", actionDirective.Name)
				}
				_resourceNameStr, ok := _resourceName.(string)
				if !ok {
					return fmt.Errorf("failed to cast Resource argument value to string in %s directive", actionDirective.Name)
				}
				if _resourceNameStr != resourceNameStr {
					continue
				}

				action := make(map[string]any)

				for _, arg := range actionDirective.Arguments {
					val, err := arg.Value.Value(nil)
					if err != nil {
						return fmt.Errorf("failed to get value of %s argument in %s directive", arg.Name, actionDirective.Name)
					}
					action[arg.Name] = val
				}
				actions = append(actions, action)
			}
		}

		resourceValue[resourceNameStr] = singleResourceValue
	}

	(*introValue)[resourceQueryField.Name] = resourceValue

	if len(resourceMapType.Fields) == 0 {
		return nil
	}

	introspectionType := document.Definitions.ForName(IntrospectionTypeName)
	if introspectionType == nil {
		return fmt.Errorf("failed to find XgenIntrospection type")
	}
	introspectionType.Fields = common.AppendFieldIfNotExists(introspectionType.Fields, resourceQueryField)

	query := document.Extensions.ForName("Query")
	if query == nil {
		return fmt.Errorf("failed to find Query type")
	}

	document.Definitions = generator.AppendDefinitionsIfNotExists(
		document.Definitions,
		resourceMapType,
		resourceDefinitionType,
		resourcePropertyType,
		actionType,
	)

	return nil

}
