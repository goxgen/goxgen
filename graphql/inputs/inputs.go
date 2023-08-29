package inputs

import "github.com/vektah/gqlparser/v2/ast"

var All = []*ast.Definition{
	PaginationInput,
	CursorPaginationInput,
	ResourceDbConfigInput,
	ResourceFieldDbConfigInput,
}
