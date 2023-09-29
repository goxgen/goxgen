package common

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	ToObjectType = directives.InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ToObjectType,
			Description: `This directive is used to define the object type`,
			Position:    pos,
			Locations: []ast.DirectiveLocation{
				ast.LocationArgumentDefinition,
				ast.LocationInputFieldDefinition,
				ast.LocationFieldDefinition,
			},
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "type",
					Type: ast.NonNullNamedType("String", nil),
				},
			},
		},
	}
)
