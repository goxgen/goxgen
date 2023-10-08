package resource

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/generator"
	"github.com/goxgen/goxgen/graphql/pagination"
	"github.com/goxgen/goxgen/graphql/sort"
	"github.com/goxgen/goxgen/runtime/gorm_initial/generated"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func SchemaGeneratorHook(schema *ast.Schema) generator.SchemaHook {
	return func(document *ast.SchemaDocument) error {
		objects := common.GetDefinedObjects(schema)
		query := &ast.Definition{
			Kind:   ast.Object,
			Name:   "Query",
			Fields: []*ast.FieldDefinition{},
		}
		mutation := &ast.Definition{
			Kind:   ast.Object,
			Name:   "Mutation",
			Fields: []*ast.FieldDefinition{},
		}

		for _, object := range objects {

			// prepare resource related types
			resourceDirective := object.Directives.ForName(consts.SchemaDefDirectiveResourceName)
			if resourceDirective != nil {
				resourceName := resourceDirective.Arguments.ForName(consts.SchemaDefResourceDirectiveArgName).Value.Raw

				sortableFieldsEnum := prepareResourceSortableFieldsEnum(object)
				document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, sortableFieldsEnum)

				sortSingleInput, sortInputObject := sort.GenerateResourceSortInputObject(resourceName)
				document.Definitions = generator.AppendDefinitionsIfNotExists(document.Definitions, sortSingleInput, sortInputObject)
			}

			// prepare resource actions related types and schema fields
			resourceActionDirectives := append(
				object.Directives.ForNames(consts.SchemaDefDirectiveActionName),
				object.Directives.ForNames(consts.SchemaDefDirectiveListActionName)...,
			)
			for _, directive := range resourceActionDirectives {
				err := prepareSchemaField(schema, query, mutation, object, directive)
				if err != nil {
					return fmt.Errorf("failed to prepare field: %w", err)
				}
			}
		}

		if len(query.Fields) > 0 {
			document.Extensions = append(document.Extensions, query)
		}

		if len(mutation.Fields) > 0 {
			document.Extensions = append(document.Extensions, mutation)
		}

		return nil
	}
}

func prepareResourceSortableFieldsEnum(object *ast.Definition) *ast.Definition {

	resourceDirective := object.Directives.ForName(consts.SchemaDefDirectiveResourceName)
	if resourceDirective == nil {
		panic("resource directive is required")
	}

	resourceName := resourceDirective.Arguments.ForName(consts.SchemaDefResourceDirectiveArgName).Value.Raw

	enumName := strings.ToUpper(resourceName + consts.ResourceSortableFieldEnumSuffix)

	enum := &ast.Definition{
		Kind:       ast.Enum,
		Name:       enumName,
		EnumValues: []*ast.EnumValueDefinition{},
	}
	for _, field := range object.Fields {
		enum.EnumValues = append(enum.EnumValues, &ast.EnumValueDefinition{
			Name: field.Name,
		})
	}
	return enum
}

func prepareSchemaField(
	schema *ast.Schema,
	query *ast.Definition,
	mutation *ast.Definition,
	object *ast.Definition,
	directive *ast.Directive,
) (err error) {
	var def *ast.Definition
	var fieldName string

	if common.IsQueryAction(directive) {
		fieldName = "where"
		def = query
	} else if common.IsMutationAction(directive) {
		fieldName = "input"
		def = mutation
	} else {
		return fmt.Errorf("failed to prepare schema field: unknown action type")
	}

	returnType, err := common.GetResourceDirectiveSingularType(schema, directive)
	if err != nil {
		return fmt.Errorf("failed to get object singular type by resource name: %w", err)
	}

	args := ast.ArgumentDefinitionList{
		{
			Name: fieldName,
			Type: &ast.Type{
				NamedType: object.Name,
			},
		},
	}

	listActionDirective := object.Directives.ForName(consts.SchemaDefDirectiveListActionName)
	if listActionDirective != nil {

		resourceName := listActionDirective.Arguments.ForName(consts.SchemaDefActionDirectiveArgResource).Value.Raw
		resourceListConfig := &generated.ListAction{}
		err = common.ArgsToStruct(listActionDirective.Arguments, resourceListConfig)
		if err != nil {
			return fmt.Errorf("failed to get resource list config: %w", err)
		}
		returnType = ast.NonNullListType(returnType, nil)
		if utils.PBool(resourceListConfig.Pagination) {
			args = append(args, &ast.ArgumentDefinition{
				Name: "pagination",
				Type: ast.NamedType(pagination.Input.Name, nil),
			})
		}

		if resourceListConfig.Sort != nil && (resourceListConfig.Sort.Disabled == nil || !utils.PBool(resourceListConfig.Sort.Disabled)) {
			args = append(args, sort.GenerateResourceQueryArgumentDefinition(resourceName))
		}
	}

	schemaQueryFieldNameArg := directive.Arguments.ForName(consts.SchemaDefActionDirectiveArgSchemaFieldName)
	if schemaQueryFieldNameArg == nil {
		return fmt.Errorf("failed to prepare schema field: %s argument is required", consts.SchemaDefActionDirectiveArgSchemaFieldName)
	}
	schemaFieldName, err := schemaQueryFieldNameArg.Value.Value(nil)
	if err != nil {
		return fmt.Errorf("failed to prepare schema field: %w", err)
	}

	schemaFieldNameStr, ok := schemaFieldName.(string)
	if !ok {
		return fmt.Errorf("failed to prepare schema field: %s argument must be a string", consts.SchemaDefActionDirectiveArgSchemaFieldName)
	}
	schemaField := &ast.FieldDefinition{
		Name:      schemaFieldNameStr,
		Arguments: args,
		Type:      returnType,
	}

	def.Fields = common.AppendFieldIfNotExists(def.Fields, schemaField)

	return nil
}
