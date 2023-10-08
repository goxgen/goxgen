package singular

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/resource/schema/actions"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	// ActionType is the enum type for the XgenResourceActionType enum
	// Values are always in upper case and snake case and should be ended with _QUERY or _MUTATION
	// depending on the type of the action.
	// You can add more values to this enum by extending the enum definition in your schema.graphql file
	// Example:
	// ```go
	// extend enum XgenResourceActionType {
	//   MY_CUSTOM_ACTION_MUTATION
	// }
	// ```
	ActionType = &ast.Definition{
		Kind: ast.Enum,
		Name: consts.SchemaDefResourceActionType,
		EnumValues: []*ast.EnumValueDefinition{
			{
				Name: consts.ActionTypeCreateMutation,
			},
			{
				Name: consts.ActionTypeReadQuery,
			},
			{
				Name: consts.ActionTypeUpdateMutation,
			},
			{
				Name: consts.ActionTypeDeleteMutation,
			},
		},
	}
	ActionDirective = directives.InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.SchemaDefDirectiveActionName,
			Description: `This directive is used to mark the object as a resource action`,
			Position:    &ast.Position{Src: &ast.Source{BuiltIn: false}},
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: consts.SchemaDefActionDirectiveArgResource,
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: consts.SchemaDefActionDirectiveActionArgAction,
					Type: ast.NonNullNamedType(ActionType.Name, nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: consts.SchemaDefActionDirectiveArgSchemaFieldName,
					Type: ast.NamedType("String", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputObject,
			},
			IsRepeatable: true,
		},
		Validate: func(directive *ast.Directive, object *ast.Definition) error {
			err := actions.PrepareActionDefaults(directive)
			if err != nil {
				return fmt.Errorf("failed to prepare action defaults: %w", err)
			}
			return nil
		},
	}
)
