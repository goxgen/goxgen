package resource

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/graphql/common"
	"github.com/goxgen/goxgen/graphql/directives"
	"github.com/goxgen/goxgen/graphql/sort"
	"github.com/goxgen/goxgen/runtime/simple_initial/generated"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	ListActionDirective = directives.InputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        consts.ListActionDirectiveName,
			Description: `This directive is used to mark the object as a resource list action`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Resource",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Action",
					Type: ast.NonNullNamedType(ListActionType.Name, nil),
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
					Name: consts.ResourceSchemaFieldName,
					Type: ast.NamedType("String", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputObject,
			},
			IsRepeatable: true,
		},
		Validate: func(directive *ast.Directive, def *ast.Definition) error {
			listActionDirective := def.Directives.ForName(consts.ListActionDirectiveName)
			if listActionDirective == nil {
				return fmt.Errorf("directive %s not found", consts.ListActionDirectiveName)
			}
			config := &generated.ListAction{}
			err := common.ArgsToStruct(listActionDirective.Arguments, config)
			//config, err := GetResourceListConfig(def)
			if err != nil {
				return err
			}
			if config.Action == ListActionTypeBrowseQuery {
				idField := def.Fields.ForName("id")
				if idField == nil {
					return fmt.Errorf("id field required for %s action", ListActionTypeBrowseQuery)
				}
			}

			err = prepareActionDefaults(directive)
			if err != nil {
				return fmt.Errorf("failed to prepare action defaults: %w", err)
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

func GetResourceListConfig(def *ast.Definition) (*XgenResourceListActionStruct, error) {
	if def == nil {
		return nil, fmt.Errorf("definition is nil")
	}
	directive := def.Directives.ForName(consts.ListActionDirectiveName)
	if directive == nil {
		return nil, fmt.Errorf("directive %s not found", consts.ListActionDirectiveName)
	}
	resource, err := directive.Arguments.ForName("Resource").Value.Value(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource argument: %w", err)
	}
	action, err := directive.Arguments.ForName("Action").Value.Value(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get action argument: %w", err)
	}
	routeArg := directive.Arguments.ForName("Route")
	var route *string
	if routeArg != nil {
		_route, _ := routeArg.Value.Value(nil)
		route = utils.StringP(_route.(string))
	}

	paginationArg := directive.Arguments.ForName("Pagination")
	var pagination *bool
	if paginationArg != nil {
		_pagination, _ := paginationArg.Value.Value(nil)
		pagination = utils.BoolP(_pagination.(bool))
	}

	return &XgenResourceListActionStruct{
		Resource:   resource.(string),
		Action:     action.(string),
		Route:      route,
		Pagination: pagination,
	}, nil
}

func IsResourceListPaginationEnabled(def *ast.Definition) (bool, error) {
	if def == nil {
		return false, fmt.Errorf("definition is nil")
	}
	directive := def.Directives.ForName(consts.ListActionDirectiveName)
	if directive == nil {
		return false, nil
	}

	paginationArg, err := directive.Arguments.ForName("Pagination").Value.Value(nil)
	if err != nil {
		return false, fmt.Errorf("failed to get pagination argument: %w", err)
	}
	paginationBool, ok := paginationArg.(bool)
	if !ok {
		return false, fmt.Errorf("type of pagination is not bool")
	}

	return paginationBool, nil
}
