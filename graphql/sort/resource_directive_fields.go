package sort

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

var ResourceListActionArgumentDefinition = &ast.ArgumentDefinition{
	Name: "Sort",
	Type: ast.NamedType(ResourceConfigInputObject.Name, nil),
	Directives: ast.DirectiveList{
		{
			Name: consts.ToObjectType,
			Arguments: ast.ArgumentList{
				{
					Name: "type",
					Value: &ast.Value{
						Kind: ast.StringValue,
						Raw:  ResourceConfigObject.Name,
					},
				},
			},
		},
	},
}
