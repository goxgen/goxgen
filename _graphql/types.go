package _graphql

import "github.com/vektah/gqlparser/v2/ast"

var Icon = &ast.Type{
	NamedType: "XgenIcon",
	Elem:      &ast.Type{NamedType: "String"},
	NonNull:   false,
	Position:  nil,
}
