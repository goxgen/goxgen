package directives

import (
	"github.com/goxgen/goxgen/graphql/enum"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	xgenResourceAction = XgenInputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        XgenResourceActionDirectiveName,
			Description: `This directive is used to mark the object as a resource action`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Resource",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Action",
					Type: ast.NonNullNamedType(enum.XgenResourceActionType.Name, nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputObject,
			},
		},
	}
)
