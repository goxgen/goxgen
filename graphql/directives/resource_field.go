package directives

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	resourceFieldDirective = FieldDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.FieldDirectiveName,
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
					Name: "DB",
					Type: ast.NamedType("XgenResourceFieldDbConfigInput", nil),
					Directives: ast.DirectiveList{
						{Name: consts.ExcludeArgumentFromType},
					},
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
			},
		},
	}
)
