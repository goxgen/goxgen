package directives

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	resourceActionFieldDirective = InputFieldDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ResourceActionFieldDirectiveName,
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
				ast.LocationInputFieldDefinition,
			},
		},
	}
)
