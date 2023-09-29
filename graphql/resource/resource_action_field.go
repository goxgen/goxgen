package resource

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

var (
	ActionFieldDirective = directives.InputFieldDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ActionFieldDirectiveName,
			Description: `This directive is used to mark the object as a resource field`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Label",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: "Description",
					Type: ast.NamedType("String", nil),
				},
				{
					Name:        consts.ActionFieldMapToArgName,
					Description: "Map field to resource field, {resource}.{field}, eg. user.id",
					Type:        ast.ListType(ast.NonNullNamedType("String", nil), nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputFieldDefinition,
			},
		},
		Validate: func(directive *ast.Directive, field *ast.FieldDefinition) error {
			resMapArg := directive.Arguments.ForName(consts.ActionFieldMapToArgName)
			if resMapArg != nil {
				resMapArgValue, err := resMapArg.Value.Value(nil)
				if err != nil {
					return fmt.Errorf("failed to get ResourceMap value: %w", err)
				}
				resMap, ok := resMapArgValue.([]any)
				if !ok {
					return fmt.Errorf("invalid ResourceMap value: %v", resMapArgValue)
				}
				for _, resMapItem := range resMap {
					resMapItemStr, ok := resMapItem.(string)
					if !ok {
						return fmt.Errorf("invalid ResourceMap item value: %v", resMapItem)
					}
					parts := strings.Split(resMapItemStr, ".")
					if len(parts) != 2 {
						return fmt.Errorf("invalid ResourceMap value: %v, should be {resource}.{field}", resMapItem)
					}
				}
			}
			return nil
		},
	}
)
