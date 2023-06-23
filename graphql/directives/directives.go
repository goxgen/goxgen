package directives

import (
	"github.com/vektah/gqlparser/v2/ast"
)

type XgenObjectDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validator  func(directive *ast.Directive, object *ast.Definition) error
}

type XgenInputObjectDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validate   func(directive *ast.Directive, object *ast.Definition) error
}

type XgenFieldDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validator  func(directive *ast.Directive, field *ast.Field) error
}
type XgenInputFieldDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validator  func(directive *ast.Directive, field *ast.Field) error
}

type XgenDirectiveDefinitionList struct {
	Object      []*XgenObjectDirectiveDefinition
	InputObject []*XgenInputObjectDirectiveDefinition
	Field       []*XgenFieldDirectiveDefinition
	InputField  []*XgenInputFieldDirectiveDefinition
}

func (xddl *XgenDirectiveDefinitionList) GetObjectDirectiveDefinition(name string) *XgenObjectDirectiveDefinition {
	for _, xdd := range xddl.Object {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (xddl *XgenDirectiveDefinitionList) GetInputObjectDirectiveDefinition(name string) *XgenInputObjectDirectiveDefinition {
	for _, xdd := range xddl.InputObject {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (xddl *XgenDirectiveDefinitionList) GetFieldDirectiveDefinition(name string) *XgenFieldDirectiveDefinition {
	for _, xdd := range xddl.Field {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (xddl *XgenDirectiveDefinitionList) GetInputFieldDirectiveDefinition(name string) *XgenInputFieldDirectiveDefinition {
	for _, xdd := range xddl.InputField {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

const (
	XgenResourceDirectiveName            = "XgenResource"
	XgenResourceActionDirectiveName      = "XgenResourceAction"
	XgenResourceFieldDirectiveName       = "XgenResourceField"
	XgenResourceListActionDirectiveName  = "XgenResourceListAction"
	XgenResourceActionFieldDirectiveName = "XgenResourceActionField"
)

func (xddl *XgenDirectiveDefinitionList) DirectiveDefinitionList() ast.DirectiveDefinitionList {
	ddl := ast.DirectiveDefinitionList{}
	for _, xdd := range xddl.Object {
		ddl = append(ddl, xdd.Definition)
	}
	for _, xdd := range xddl.InputObject {
		ddl = append(ddl, xdd.Definition)
	}
	for _, xdd := range xddl.Field {
		ddl = append(ddl, xdd.Definition)
	}
	for _, xdd := range xddl.InputField {
		ddl = append(ddl, xdd.Definition)
	}
	return ddl
}

var (
	pos = &ast.Position{Src: &ast.Source{BuiltIn: false}}

	All = &XgenDirectiveDefinitionList{
		Object: []*XgenObjectDirectiveDefinition{
			&xgenResource,
		},
		InputObject: []*XgenInputObjectDirectiveDefinition{
			&xgenResourceAction,
			&xgenResourceListAction,
		},
		Field: []*XgenFieldDirectiveDefinition{
			&xgenResourceField,
		},
		InputField: []*XgenInputFieldDirectiveDefinition{
			&xgenResourceActionField,
		},
	}
)

func mergeDirectiveDefs(directive ast.DirectiveDefinition, new ast.DirectiveDefinition) *ast.DirectiveDefinition {
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

	return &directive
}
