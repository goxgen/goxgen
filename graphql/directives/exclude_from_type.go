package directives

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	excludeArgumentFromTypeDirective = InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ExcludeArgumentFromType,
			Description: `This directive is used to exclude the argument from the type`,
			Position:    pos,
			Locations: []ast.DirectiveLocation{
				ast.LocationArgumentDefinition,
			},
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "exclude",
					Type: ast.NamedType("Boolean", nil),
				},
			},
		},
	}
)
