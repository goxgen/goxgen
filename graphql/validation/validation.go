package validation

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/vektah/gqlparser/v2/ast"
)

func SchemaGeneratorHook(schema *ast.Schema) generator.SchemaHook {
	return func(_ *ast.SchemaDocument) error {
		objects := common.GetDefinedObjects(schema)
		for _, object := range objects {
			resActionDirs := append(
				object.Directives.ForNames(consts.ActionDirectiveName),
				object.Directives.ForNames(consts.ListActionDirectiveName)...,
			)
			for _, resActionDir := range resActionDirs {
				xgenDirDef := directives.Bundle.GetInputObjectDirectiveDefinition(resActionDir.Name)
				if xgenDirDef != nil && xgenDirDef.Validate != nil {
					err := xgenDirDef.Validate(resActionDir, object)
					if err != nil {
						return fmt.Errorf("failed to validate object %s: %w", object.Name, err)
					}
				}
			}

			for _, field := range object.Fields {
				resActionFieldDirs := field.Directives.ForNames(consts.ActionFieldDirectiveName)
				for _, resActionFieldDir := range resActionFieldDirs {
					xgenDirDef := directives.Bundle.GetInputFieldDirectiveDefinition(resActionFieldDir.Name)
					if xgenDirDef != nil && xgenDirDef.Validate != nil {
						err := xgenDirDef.Validate(resActionFieldDir, field)
						if err != nil {
							return fmt.Errorf("failed to validate field %s: %w", field.Name, err)
						}
					}
				}
			}
		}
		return nil
	}
}
