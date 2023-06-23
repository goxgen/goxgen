package directives

import (
	"github.com/goxgen/goxgen/graphql/enum"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	pos         = &ast.Position{Src: &ast.Source{BuiltIn: false}}
	XgenVersion = ast.DirectiveDefinition{
		Name:        "XgenVersion",
		Description: `This directive is used to specify the version of the schema.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "version",
				Type: ast.NonNullNamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationSchema,
		},
	}
	XgenResource = ast.DirectiveDefinition{
		Name:        "XgenResource",
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
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
		},
	}
	XgenResourceAction = ast.DirectiveDefinition{
		Name:        "XgenResourceAction",
		Description: `This directive is used to mark the object as a resource action`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Resource",
				Type: ast.NonNullNamedType("String", nil),
			},
			{
				Name: "Action",
				Type: ast.NonNullNamedType(enum.XgenResourceActionType.Name, nil),
			},
			{
				Name: "Route",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationInputObject,
		},
	}
	XgenResourceField = ast.DirectiveDefinition{
		Name:        "XgenResourceField",
		Description: `This directive is used to mark the object as a resource field`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Label",
				Type: ast.NamedType("String", nil),
			},
			{
				Name: "Description",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationFieldDefinition,
		},
	}
	XgenResourceActionField = mergeDirectiveDefs(XgenResourceField, ast.DirectiveDefinition{
		Name:      "XgenResourceActionField",
		Locations: []ast.DirectiveLocation{ast.LocationInputFieldDefinition},
	})

	All = []*ast.DirectiveDefinition{
		&XgenVersion,

		// Resource related directives
		&XgenResource,
		&XgenResourceAction,
		&XgenResourceListAction,
		&XgenResourceField,
		&XgenResourceActionField,
	}
)

func mergeDirectiveDefs(directive ast.DirectiveDefinition, new ast.DirectiveDefinition) ast.DirectiveDefinition {
	directive.Name = new.Name

	if len(new.Locations) != 0 {
		directive.Locations = new.Locations
	}
	if new.Position != nil {
		directive.Position = new.Position
	}
	if len(new.Arguments) != 0 {
		directive.Arguments = new.Arguments
	}
	if new.Description != "" {
		directive.Description = new.Description
	}

	return directive
}
