package directives

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func prepareActionDefaults(dir *ast.Directive) error {
	schemaQueryFieldNameArg := dir.Arguments.ForName(consts.ResourceSchemaFieldName)
	if schemaQueryFieldNameArg == nil {
		resourceArg, err := dir.Arguments.ForName("Resource").Value.Value(nil)
		if err != nil {
			return fmt.Errorf("failed to get resource argument: %w", err)
		}
		resourceArgStr, ok := resourceArg.(string)
		if !ok {
			return fmt.Errorf("resource argument is not string")
		}

		resActionEnum := dir.Arguments.ForName("Action").Value.Raw
		resAction := strings.TrimSuffix(resActionEnum, "_QUERY")
		resAction = strings.TrimSuffix(resAction, "_MUTATION")
		resAction = strings.ToLower(resAction)

		dir.Arguments = append(dir.Arguments, &ast.Argument{
			Name: consts.ResourceSchemaFieldName,
			Value: &ast.Value{
				Kind: ast.StringValue,
				Raw:  resourceArgStr + "_" + resAction,
			},
		})
	}
	return nil
}
