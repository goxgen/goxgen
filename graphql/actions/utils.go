package actions

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

// PrepareActionDefaults prepares default values for action directive
func PrepareActionDefaults(dir *ast.Directive) error {
	schemaQueryFieldNameArg := dir.Arguments.ForName(consts.SchemaDefActionDirectiveArgSchemaFieldName)
	if schemaQueryFieldNameArg == nil {
		resourceArg, err := dir.Arguments.ForName(consts.SchemaDefActionDirectiveArgResource).Value.Value(nil)
		if err != nil {
			return fmt.Errorf("failed to get resource argument: %w", err)
		}
		resourceArgStr, ok := resourceArg.(string)
		if !ok {
			return fmt.Errorf("resource argument is not string")
		}

		resActionEnum := dir.Arguments.ForName(consts.SchemaDefActionDirectiveActionArgAction).Value.Raw
		resAction := strings.TrimSuffix(resActionEnum, "_QUERY")
		resAction = strings.TrimSuffix(resAction, "_MUTATION")
		resAction = strings.ToLower(resAction)

		dir.Arguments = append(dir.Arguments, &ast.Argument{
			Name: consts.SchemaDefActionDirectiveArgSchemaFieldName,
			Value: &ast.Value{
				Kind: ast.StringValue,
				Raw:  resourceArgStr + "_" + resAction,
			},
		})
	}
	return nil
}

func PrepareListActionDefaults(dir *ast.Directive) error {
	sortArg := dir.Arguments.ForName(consts.SortQueryArgumentName)
	if sortArg == nil {
		dir.Arguments = append(dir.Arguments, &ast.Argument{
			Name: consts.SortQueryArgumentName,
			Value: &ast.Value{
				Kind: ast.ObjectValue,
				Raw:  "{}",
			},
		})
	}
	return nil
}
