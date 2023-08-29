package validation

import (
	"fmt"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

func SchemaGeneratorHook(schema *ast.Schema) generator.SchemaHook {
	return func(_ *ast.SchemaDocument) error {
		objects := common.GetDefinedObjects(schema)
		for _, object := range objects {
			resourceActionDirectives := directives.GetResourceActionDirectives(object)
			for _, directive := range resourceActionDirectives {
				xgenDirDef := directives.Bundle.GetInputObjectDirectiveDefinition(directive.Name)
				if xgenDirDef != nil && xgenDirDef.Validate != nil {
					err := xgenDirDef.Validate(directive, object)
					if err != nil {
						return fmt.Errorf("failed to validate object %s: %w", object.Name, err)
					}
				}
			}

		}
		return nil
	}
}
