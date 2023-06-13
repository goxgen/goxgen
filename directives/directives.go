package directives

import "github.com/vektah/gqlparser/v2/ast"

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
	XgenObjectAuth = ast.DirectiveDefinition{
		Name:        "XgenObjectAuth",
		Description: `This directive is used to specify the authorization of the object.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Resource",
				Type: ast.NamedType("String", nil),
			},
			{
				Name: "Scope",
				Type: ast.NamedType("String", nil),
			},
			{
				Name: "Roles",
				Type: ast.ListType(ast.NonNullNamedType("String", nil), nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenFieldAuth = mergeDirectives(
		XgenObjectAuth,
		ast.DirectiveDefinition{
			Name:        "XgenFieldAuth",
			Description: `This directive is used to specify the authorization of the field.`,
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	)
	XgenObjectIcon = ast.DirectiveDefinition{
		Name:        "XgenObjectIcon",
		Description: `This directive is used to specify the icon of the object.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Icon",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenFieldIcon = mergeDirectives(
		XgenObjectIcon,
		ast.DirectiveDefinition{
			Name:        "XgenFieldIcon",
			Description: `This directive is used to specify the icon of the field.`,
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	)
	XgenObjectLabel = ast.DirectiveDefinition{
		Name:        "XgenObjectLabel",
		Description: `This directive is used to specify the label of the object.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Label",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenFieldLabel = mergeDirectives(
		XgenObjectLabel,
		ast.DirectiveDefinition{
			Name:        "XgenFieldLabel",
			Description: `This directive is used to specify the label of the field.`,
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	)
	XgenObjectDescription = ast.DirectiveDefinition{
		Name:        "XgenObjectDescription",
		Description: `This directive is used to specify the description of the object.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Description",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenFieldDescription = mergeDirectives(
		XgenObjectDescription,
		ast.DirectiveDefinition{
			Name:        "XgenFieldDescription",
			Description: `This directive is used to specify the description of the field.`,
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	)
	XgenObjectRoute = ast.DirectiveDefinition{
		Name:        "XgenObjectRoute",
		Description: `This directive is used to specify the route of the object.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Route",
				Type: ast.NamedType("String", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenObjectIndexRoute = ast.DirectiveDefinition{
		Name:        "XgenObjectIndexRoute",
		Description: `This directive is used to mark the object as the index route.`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "UseAsIndexRoute",
				Type: ast.NamedType("Boolean", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationObject,
			ast.LocationInputObject,
		},
	}
	XgenFieldCode = ast.DirectiveDefinition{
		Name:        "XgenFieldCode",
		Description: `This directive is used to mark the field as a code field`,
		Position:    pos,
		Arguments: ast.ArgumentDefinitionList{
			{
				Name: "Language",
				Type: ast.NamedType("String", nil),
			},
			{
				Name: "TabSize",
				Type: ast.NamedType("Int", nil),
			},
		},
		Locations: []ast.DirectiveLocation{
			ast.LocationFieldDefinition,
			ast.LocationInputFieldDefinition,
		},
	}

	All = []*ast.DirectiveDefinition{
		&XgenVersion,
		&XgenObjectAuth,
		&XgenFieldAuth,
		&XgenObjectIcon,
		&XgenFieldIcon,
		&XgenObjectLabel,
		&XgenFieldLabel,
		&XgenObjectDescription,
		&XgenFieldDescription,
		&XgenObjectRoute,
		&XgenObjectIndexRoute,
		&XgenFieldCode,
	}
)

func mergeDirectives(directive ast.DirectiveDefinition, new ast.DirectiveDefinition) ast.DirectiveDefinition {
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

	return directive
}
