package resource

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	ActionDirective = directives.InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ActionDirectiveName,
			Description: `This directive is used to mark the object as a resource action`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Resource",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Action",
					Type: ast.NonNullNamedType(ActionType.Name, nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: consts.ResourceSchemaFieldName,
					Type: ast.NamedType("String", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputObject,
			},
			IsRepeatable: true,
		},
		Validate: func(directive *ast.Directive, object *ast.Definition) error {
			err := prepareActionDefaults(directive)
			if err != nil {
				return fmt.Errorf("failed to prepare action defaults: %w", err)
			}
			return nil
		},
	}
)
