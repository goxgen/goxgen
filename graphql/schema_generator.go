package graphql

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"os"
	"strings"
)

// SchemaHook is a function that can be used to modify the schema before it is
type SchemaHook func(*ast.Schema) error

// OutputWriter is a function that can be used to write the schema to a custom output.
type OutputWriter func(*ast.Schema) error

// SchemaGenerator is a struct that can be used to generate a schema from a list of schema hooks.
// The schema can then be written to a file or a custom output.
type SchemaGenerator struct {
	Path         string
	OutputWriter OutputWriter
	SchemaHooks  []SchemaHook
	Header       string
	Footer       string
}

// printSchema prints the given schema to a string.
func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb, formatter.WithIndent("  ")).
		FormatSchema(schema)
	return sb.String()
}

// buildSchema builds a schema from the given schema hooks.
func (e *SchemaGenerator) buildSchema() (s *ast.Schema, err error) {
	s = &ast.Schema{
		Directives:    make(map[string]*ast.DirectiveDefinition),
		Types:         make(map[string]*ast.Definition),
		PossibleTypes: make(map[string][]*ast.Definition),
		Implements:    make(map[string][]*ast.Definition),
	}

	for _, h := range e.SchemaHooks {
		if err = h(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// GenerateOutput generates the schema and writes it to the given path.
func (e *SchemaGenerator) GenerateOutput() error {
	schema, err := e.buildSchema()
	if err != nil {
		return err
	}
	if e.OutputWriter == nil {
		if e.Path == "" {
			return nil
		}
		content := strings.Join(
			[]string{
				e.Header,
				printSchema(schema),
				e.Footer,
			},
			"\n",
		)
		content = strings.Trim(content, "\n")

		return os.WriteFile(
			e.Path,
			[]byte(content),
			0644,
		)
	}
	return err
}
