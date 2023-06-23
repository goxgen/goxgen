package directives

import (
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
	"testing"
)

// test mergeDirectiveDefs function with testify/assert
func TestMergeDirectiveDefs(t *testing.T) {
	def1 := ast.DirectiveDefinition{
		Name: "Def1",
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Arg1",
				Type: ast.NonNullNamedType("String", nil),
			},
			{
				Name: "Arg2",
				Type: ast.NonNullNamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
		},
	}

	def2 := ast.DirectiveDefinition{
		Name:        "New Name",
		Description: "New Desc",
		Locations: []ast.DirectiveLocation{
			ast.LocationFieldDefinition,
		},
	}

	merged := mergeDirectiveDefs(def1, def2)

	_assert := assert.New(t)
	_assert.Equal(2, len(merged.Arguments))
	_assert.Equal("New Name", merged.Name)
	_assert.Equal("New Desc", merged.Description)
	_assert.Equal(
		[]ast.DirectiveLocation{
			ast.LocationFieldDefinition,
		},
		merged.Locations,
	)
}
