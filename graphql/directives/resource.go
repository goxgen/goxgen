package directives

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	resourceDirective = ObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ResourceDirectiveName,
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
				{
					Name: "DB",
					Type: ast.NamedType("XgenResourceDbConfigInput", nil),
					Directives: ast.DirectiveList{
						{Name: consts.ExcludeArgumentFromType},
					},
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationObject,
			},
		},
	}
)
