package field

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	FieldDirective = directives.FieldDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.SchemaDefDirectiveFieldName,
			Description: `This directive is used to mark the object as a resource field`,
			Position:    &ast.Position{Src: &ast.Source{BuiltIn: false}},
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: consts.SchemaDefFieldDirectiveArgLabel,
					Type: ast.NamedType("String", nil),
				},
				{
					Name: consts.SchemaDefFieldDirectiveArgDescription,
					Type: ast.NamedType("String", nil),
				},
				{
					Name: consts.SchemaDefFieldDirectiveArgDb,
					Type: ast.NamedType(consts.SchemaDefFieldDbConfigInputType, nil),
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
