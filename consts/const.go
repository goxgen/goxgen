package consts

const (
	GeneratedGqlgenPackageName = "generated"
	GeneratedModelsPackageName = "models"

	ExcludeArgumentFromType = "ExcludeArgumentFromType"
	ToObjectType            = "ToObjectType"

	ResourceSortableFieldEnumSuffix = "_SORTABLE_FIELD"
	SingleSortInputSuffix           = "SingleSortInput"
	SortInputSuffix                 = "SortInput"
	SortQueryArgumentName           = "sort"

	MapToGolangStructTagName = "mapto"

	// Field names
	SchemaDefActionDirectiveArgSchemaFieldName = "SchemaFieldName"
	SchemaDefActionDirectiveArgResource        = "Resource"
	SchemaDefActionDirectiveActionArgAction    = "Action"

	SortDirectionAsc  = "ASC"
	SortDirectionDesc = "DESC"

	// Schema definition constants
	SchemaDefResourceActionType     = "XgenResourceActionType"
	SchemaDefResourceListActionType = "XgenResourceListActionType"
	SchemaDefFieldDbConfigInputType = "XgenResourceFieldDbConfigInput"

	SchemaDefDirectiveResourceName    = "Resource"
	SchemaDefDirectiveActionName      = "Action"
	SchemaDefDirectiveListActionName  = "ListAction"
	SchemaDefDirectiveActionFieldName = "ActionField"
	SchemaDefDirectiveFieldName       = "Field"

	SchemaDefFieldDirectiveArgLabel       = "Label"
	SchemaDefFieldDirectiveArgDescription = "Description"
	SchemaDefFieldDirectiveArgDb          = "DB"
	SchemaDefResourceDirectiveArgName     = "Name"
	SchemaDefResourceDirectiveArgDb       = "DB"
	SchemaDefActionFieldDirectiveArgMapTo = "MapTo"

	ActionTypeCreateMutation                   = "CREATE_MUTATION"
	ActionTypeReadQuery                        = "READ_QUERY"
	ActionTypeUpdateMutation                   = "UPDATE_MUTATION"
	ActionTypeDeleteMutation                   = "DELETE_MUTATION"
	SchemaDefListActionTypeBrowseQuery         = "BROWSE_QUERY"
	SchemaDefListActionTypeBatchDeleteMutation = "BATCH_DELETE_MUTATION"
)
