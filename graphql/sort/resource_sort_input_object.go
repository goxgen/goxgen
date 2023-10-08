package sort

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func ResourceSingleSortInputObjectName(resourceName string) string {
	return strcase.ToCamel(resourceName + consts.SingleSortInputSuffix)
}

func ResourceSortInputObjectName(resourceName string) string {
	return strcase.ToCamel(resourceName + consts.SortInputSuffix)
}

func GenerateResourceSortInputObject(resourceName string) (singleInput *ast.Definition, input *ast.Definition) {
	enumName := strings.ToUpper(resourceName) + consts.ResourceSortableFieldEnumSuffix
	singleInput = &ast.Definition{
		Kind: ast.InputObject,
		Name: ResourceSingleSortInputObjectName(resourceName),
		Fields: ast.FieldList{
			{
				Name: "field",
				Type: ast.NonNullNamedType(enumName, nil),
			},
			{
				Name: "direction",
				Type: ast.NamedType(DirectionEnum.Name, nil),
			},
		},
		Interfaces: []string{},
	}

	input = &ast.Definition{
		Kind: ast.InputObject,
		Name: ResourceSortInputObjectName(resourceName),
		Fields: ast.FieldList{
			{
				Name: "by",
				Type: ast.ListType(
					ast.NonNullNamedType(singleInput.Name, nil),
					nil,
				),
			},
		},
	}
	return singleInput, input
}
