package directives

import "github.com/vektah/gqlparser/v2/ast"

var (
	xgenResourceField = XgenFieldDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        XgenResourceFieldDirectiveName,
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
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
			},
		},
	}
)
