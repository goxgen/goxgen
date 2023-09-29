package resource

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/db"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	ActionTypeCreateMutation          = "CREATE_MUTATION"
	ActionTypeReadQuery               = "READ_QUERY"
	ActionTypeUpdateMutation          = "UPDATE_MUTATION"
	ActionTypeDeleteMutation          = "DELETE_MUTATION"
	ListActionTypeBrowseQuery         = "BROWSE_QUERY"
	ListActionTypeBatchDeleteMutation = "BATCH_DELETE_MUTATION"
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
		Name: "XgenResourceActionType",
		EnumValues: []*ast.EnumValueDefinition{
			{
				Name: ActionTypeCreateMutation,
			},
			{
				Name: ActionTypeReadQuery,
			},
			{
				Name: ActionTypeUpdateMutation,
			},
			{
				Name: ActionTypeDeleteMutation,
			},
		},
	}
	ListActionType = &ast.Definition{
		Kind: ast.Enum,
		Name: "XgenResourceListActionType",
		EnumValues: []*ast.EnumValueDefinition{
			{
				Name: ListActionTypeBrowseQuery,
			},
			{
				Name: ListActionTypeBatchDeleteMutation,
			},
		},
	}

	AllDefinitions = []*ast.Definition{
		ActionType,
		ListActionType,
	}
)

var (
	pos       = &ast.Position{Src: &ast.Source{BuiltIn: false}}
	Directive = directives.ObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ResourceDirectiveName,
			Description: `This directive is used to mark the object as a resource`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Name",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: "Primary",
					Type: ast.NamedType("Boolean", nil),
				},
				{
					Name: "DB",
					Type: ast.NamedType(db.ResourceConfigInput.Name, nil),
					Directives: ast.DirectiveList{
						{Name: consts.ExcludeArgumentFromType},
					},
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationObject,
			},
		},
	}
)
