package sort

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

func GenerateResourceQueryArgumentDefinition(resourceName string) *ast.ArgumentDefinition {
	return &ast.ArgumentDefinition{
		Name: consts.SortQueryArgumentName,
		Type: ast.NamedType(ResourceSortInputObjectName(resourceName), nil),
	}
}
