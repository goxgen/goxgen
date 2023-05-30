package graphql

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"strings"
)

func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb, formatter.WithIndent("  ")).
		FormatSchema(schema)
	return sb.String()
}

func BuildSchema() string {
	schema := &ast.Schema{
		Description: "The Xgen schema",
		Types: map[string]*ast.Definition{
			"XgenIconType": {
				Kind:        "OBJECT",
				Description: "asd",
				Name:        "XgenIconType",
				Fields: []*ast.FieldDefinition{
					{
						Name: "name",
						Type: &ast.Type{
							NamedType: "String",
							NonNull:   true,
						},
					},
					{
						Name: "src",
						Type: &ast.Type{
							NamedType: "String",
						},
					},
					{
						Name: "cssClasses",
						Type: &ast.Type{
							NamedType: "String",
						},
					},
				},
			},
		},
		Directives: map[string]*ast.DirectiveDefinition{
			"XgenIcon": {
				Name: "XgenIcon",
				Arguments: []*ast.ArgumentDefinition{
					{
						Name: "icon",
						Type: Icon,
					},
				},
				Position: &ast.Position{
					Src: &ast.Source{
						BuiltIn: true,
					},
				},
			},
		},
	}
	return printSchema(schema)
}
