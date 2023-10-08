package directives

import (
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
)

type ObjectDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validator  func(directive *ast.Directive, object *ast.Definition) error
}

type InputObjectDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validate   func(directive *ast.Directive, object *ast.Definition) error
}

type FieldDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validate   func(directive *ast.Directive, field *ast.FieldDefinition) error
}

type InputFieldDirectiveDefinition struct {
	Definition *ast.DirectiveDefinition
	Validate   func(directive *ast.Directive, field *ast.FieldDefinition) error
}

type DirectiveDefinitionBundle struct {
	Object      []*ObjectDirectiveDefinition
	InputObject []*InputObjectDirectiveDefinition
	Field       []*FieldDirectiveDefinition
	InputField  []*InputFieldDirectiveDefinition
}

func (ddb *DirectiveDefinitionBundle) GetObjectDirectiveDefinition(name string) *ObjectDirectiveDefinition {
	for _, xdd := range ddb.Object {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (ddb *DirectiveDefinitionBundle) GetInputObjectDirectiveDefinition(name string) *InputObjectDirectiveDefinition {
	for _, xdd := range ddb.InputObject {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (ddb *DirectiveDefinitionBundle) GetFieldDirectiveDefinition(name string) *FieldDirectiveDefinition {
	for _, xdd := range ddb.Field {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (ddb *DirectiveDefinitionBundle) GetInputFieldDirectiveDefinition(name string) *InputFieldDirectiveDefinition {
	for _, xdd := range ddb.InputField {
		if xdd.Definition.Name == name {
			return xdd
		}
	}
	return nil
}

func (ddb *DirectiveDefinitionBundle) DirectiveDefinitionList() ast.DirectiveDefinitionList {
	_ddl := ast.DirectiveDefinitionList{}
	for _, xdd := range ddb.Object {
		_ddl = append(_ddl, xdd.Definition)
	}
	for _, xdd := range ddb.InputObject {
		_ddl = append(_ddl, xdd.Definition)
	}
	for _, xdd := range ddb.Field {
		_ddl = append(_ddl, xdd.Definition)
	}
	for _, xdd := range ddb.InputField {
		_ddl = append(_ddl, xdd.Definition)
	}
	return _ddl
}

func GetInputFieldDirectives(definition *ast.Definition) []*ast.Directive {
	return definition.Directives.ForNames(consts.SchemaDefDirectiveActionFieldName)
}

func GetObjectFieldDirectives(definition *ast.Definition) []*ast.Directive {
	return definition.Directives.ForNames(consts.SchemaDefDirectiveFieldName)
}

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
