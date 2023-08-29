package inputs

import "github.com/vektah/gqlparser/v2/ast"

var (
	PaginationInput = &ast.Definition{
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

	CursorPaginationInput = &ast.Definition{
		Kind: ast.InputObject,
		Name: "XgenCursorPaginationInput",
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
