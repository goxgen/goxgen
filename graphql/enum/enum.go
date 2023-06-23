package enum

import "github.com/vektah/gqlparser/v2/ast"

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
				Name: "CREATE_MUTATION",
			},
			{
				Name: "READ_QUERY",
			},
			{
				Name: "UPDATE_MUTATION",
			},
			{
				Name: "DELETE_MUTATION",
			},
		},
	}
	XgenResourceListActionType = &ast.Definition{
		Kind: ast.Enum,
		Name: "XgenResourceListActionType",
		EnumValues: []*ast.EnumValueDefinition{
			{
				Name: "LIST_QUERY",
			},
			{
				Name: "BATCH_DELETE_MUTATION",
			},
		},
	}

	All = []*ast.Definition{
		XgenResourceActionType,
		XgenResourceListActionType,
	}
)
