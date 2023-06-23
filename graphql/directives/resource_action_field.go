package directives

import "github.com/vektah/gqlparser/v2/ast"

var (
	xgenResourceActionField = XgenInputFieldDirectiveDefinition{
		Definition: mergeDirectiveDefs(*xgenResourceField.Definition, ast.DirectiveDefinition{
			Name:      XgenResourceActionFieldDirectiveName,
			Locations: []ast.DirectiveLocation{ast.LocationInputFieldDefinition},
		}),
	}
)
