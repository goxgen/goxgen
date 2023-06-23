package directives

import (
	"fmt"
	"github.com/goxgen/goxgen/graphql/enum"
	"github.com/goxgen/goxgen/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	xgenResourceListAction = XgenInputObjectDirectiveDefinition{
		Definition: &ast.DirectiveDefinition{
			Name:        XgenResourceListActionDirectiveName,
			Description: `This directive is used to mark the object as a resource list action`,
			Position:    pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "Resource",
					Type: ast.NonNullNamedType("String", nil),
				},
				{
					Name: "Action",
					Type: ast.NonNullNamedType(enum.XgenResourceListActionType.Name, nil),
				},
				{
					Name: "Route",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: "Pagination",
					Type: ast.NamedType("Boolean", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationInputObject,
			},
		},
		Validate: func(dir *ast.Directive, def *ast.Definition) error {
			config, err := GetResourceListConfig(def)
			if err != nil {
				return err
			}
			if config.Action == enum.ListActionTypeBrowseQuery {
				idField := def.Fields.ForName("id")
				if idField == nil {
					return fmt.Errorf("id field required for %s action", enum.ListActionTypeBrowseQuery)
				}
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
	directive := def.Directives.ForName(XgenResourceListActionDirectiveName)
	if directive == nil {
		return nil, fmt.Errorf("directive %s not found", XgenResourceListActionDirectiveName)
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
	directive := def.Directives.ForName(XgenResourceListActionDirectiveName)
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
