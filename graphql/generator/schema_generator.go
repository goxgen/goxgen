package generator

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"os"
	"strings"
)

// SchemaHook is a function that can be used to modify the schema before it is
type SchemaHook func(*ast.SchemaDocument) error

// OutputWriter is a function that can be used to write the schema to a custom output.
type OutputWriter func(*ast.Schema, *ast.SchemaDocument) error

// SchemaGenerator is a struct that can be used to generate a schema from a list of schema hooks.
// The schema can then be written to a file or a custom output.
type SchemaGenerator struct {
	path         string
	outputWriter OutputWriter
	schemaHooks  []SchemaHook
	header       string
	footer       string
}

func NewSchemaGenerator() *SchemaGenerator {
	return &SchemaGenerator{}
}

func (e *SchemaGenerator) WithPath(path string) *SchemaGenerator {
	e.path = path
	return e
}

func (e *SchemaGenerator) WithOutputWriter(outputWriter OutputWriter) *SchemaGenerator {
	e.outputWriter = outputWriter
	return e
}

func (e *SchemaGenerator) WithSchemaHooks(schemaHooks ...SchemaHook) *SchemaGenerator {
	e.schemaHooks = schemaHooks
	return e
}

func (e *SchemaGenerator) WithHeader(header string) *SchemaGenerator {
	e.header = header
	return e
}

func (e *SchemaGenerator) WithFooter(footer string) *SchemaGenerator {
	e.footer = footer
	return e
}

// printSchema prints the given schema to a string.
func printSchema(schemaDocument *ast.SchemaDocument) (schemaDocumentStr string) {
	schemaDocumentBuilder := &strings.Builder{}
	formatter.
		NewFormatter(schemaDocumentBuilder, formatter.WithIndent("  ")).
		FormatSchemaDocument(schemaDocument)
	return schemaDocumentBuilder.String()
}

// buildSchema builds a schema from the given schema hooks.
func (e *SchemaGenerator) buildSchema() (sd *ast.SchemaDocument, err error) {
	sd = &ast.SchemaDocument{}

	for _, h := range e.schemaHooks {
		if err = h(sd); err != nil {
			return nil, err
		}
	}
	return sd, nil
}

// GenerateOutput generates the schema and writes it to the given path.
func (e *SchemaGenerator) GenerateOutput() error {
	schemaDocument, err := e.buildSchema()
	if err != nil {
		return err
	}
	if e.outputWriter == nil {

		schemaDocumentStr := printSchema(schemaDocument)
		if e.path == "" {
			return nil
		}

		schemaContent := strings.Join(
			[]string{
				e.header,
				schemaDocumentStr,
				e.footer,
			},
			"\n",
		)
		schemaContent = strings.Trim(schemaContent, "\n")

		return os.WriteFile(
			e.path,
			[]byte(schemaContent),
			0644,
		)
	}
	return err
}

func AppendDefinitionsIfNotExists(defs ast.DefinitionList, defsToAppend ...*ast.Definition) []*ast.Definition {
	for _, def := range defsToAppend {
		if defs.ForName(def.Name) != nil {
			continue
		}
		defs = append(defs, def)
	}
	return defs
}
