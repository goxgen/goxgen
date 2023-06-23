package inputs

import "github.com/vektah/gqlparser/v2/ast"

var (
	XgenPaginationInput = &ast.Definition{
		Kind: ast.InputObject,
		Name: "XgenPaginationInput",
		Fields: ast.FieldList{
			{
				Name: "page",
				Type: ast.NonNullNamedType("Int", nil),
			},
			{
				Name: "limit",
				Type: ast.NonNullNamedType("Int", nil),
			},
		},
	}

	XgenCursorPaginationInput = &ast.Definition{
		Kind: ast.InputObject,
		Name: "XgenPaginationInput",
		Fields: ast.FieldList{
			{
				Name: "first",
				Type: ast.NonNullNamedType("Int", nil),
			},
			{
				Name: "after",
				Type: ast.NamedType("String", nil),
			},
			{
				Name: "last",
				Type: ast.NonNullNamedType("Int", nil),
			},
			{
				Name: "before",
				Type: ast.NamedType("String", nil),
			},
		},
	}
)
