package list

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/resource/schema/actions"
	"github.com/goxgen/goxgen/graphql/sort"
	"github.com/goxgen/goxgen/runtime/simple_initial/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	ActionType = &ast.Definition{
		Kind: ast.Enum,
		Name: consts.SchemaDefResourceListActionType,
		EnumValues: []*ast.EnumValueDefinition{
			{
				Name: consts.SchemaDefListActionTypeBrowseQuery,
			},
			{
				Name: consts.SchemaDefListActionTypeBatchDeleteMutation,
			},
		},
	}
	ActionDirective = directives.InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.SchemaDefDirectiveListActionName,
			Description: `This directive is used to mark the object as a resource list action`,
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
					Name: "Pagination",
					Type: ast.NamedType("Boolean", nil),
				},
				sort.ResourceListActionArgumentDefinition,
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
		Validate: func(directive *ast.Directive, def *ast.Definition) error {
			listActionDirective := def.Directives.ForName(consts.SchemaDefDirectiveListActionName)
			if listActionDirective == nil {
				return fmt.Errorf("directive %s not found", consts.SchemaDefDirectiveListActionName)
			}
			config := &generated.ListAction{}
			err := common.ArgsToStruct(listActionDirective.Arguments, config)
			//config, err := GetResourceListConfig(def)
			if err != nil {
				return err
			}
			if config.Action == consts.SchemaDefListActionTypeBrowseQuery {
				idField := def.Fields.ForName("id")
				if idField == nil {
					return fmt.Errorf("id field required for %s action", consts.SchemaDefListActionTypeBrowseQuery)
				}
			}

			err = actions.PrepareActionDefaults(directive)
			if err != nil {
				return fmt.Errorf("failed to prepare action defaults: %w", err)
			}

			err = actions.PrepareListActionDefaults(directive)
			if err != nil {
				return fmt.Errorf("failed to prepare list action defaults: %w", err)
			}

			return nil
		},
	}
)

type XgenResourceListActionStruct struct {
	Resource   string
	Action     string
	Route      *string
	Pagination *bool
}
