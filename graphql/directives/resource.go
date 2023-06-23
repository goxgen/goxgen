package directives

import "github.com/vektah/gqlparser/v2/ast"

var (
	xgenResource = XgenObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        XgenResourceDirectiveName,
			Description: `This directive is used to mark the object as a resource`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Name",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: "Primary",
					Type: ast.NamedType("Boolean", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationObject,
			},
		},
	}
)
