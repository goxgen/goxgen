package enum

import "github.com/vektah/gqlparser/v2/ast"

const (
	ActionTypeCreateMutation          = "CREATE_MUTATION"
	ActionTypeReadQuery               = "READ_QUERY"
	ActionTypeUpdateMutation          = "UPDATE_MUTATION"
	ActionTypeDeleteMutation          = "DELETE_MUTATION"
	ListActionTypeBrowseQuery         = "BROWSE_QUERY"
	ListActionTypeBatchDeleteMutation = "BATCH_DELETE_MUTATION"
)

var (
	// XgenResourceActionType is the enum type for the XgenResourceActionType enum
	// Values are always in upper case and snake case and should be ended with _QUERY or _MUTATION
	// depending on the type of the action.
	// You can add more values to this enum by extending the enum definition in your schema.graphql file
	// Example:
	// ```go
	// extend enum XgenResourceActionType {
	//   MY_CUSTOM_ACTION_MUTATION
	// }
	// ```
	XgenResourceActionType = &ast.Definition{
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
	XgenResourceListActionType = &ast.Definition{
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

	All = []*ast.Definition{
		XgenResourceActionType,
		XgenResourceListActionType,
	}
)
