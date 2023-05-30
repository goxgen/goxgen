package graphql

import "github.com/vektah/gqlparser/v2/ast"

func SetIcon(iconName string) *ast.Directive {
	return &ast.Directive{
		Name: "XgenIcon",
		Arguments: []*ast.Argument{
			{
				Name: "icon",
				Value: &ast.Value{
					Raw:          iconName,
					ExpectedType: Icon,
				},
			},
		},
	}
}
