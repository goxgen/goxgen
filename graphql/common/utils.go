package common

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

// IsXgenDirectiveDefinition checks if directive is xgen directive
func IsXgenDirectiveDefinition(directive *ast.DirectiveDefinition) bool {
	return directive.Name == consts.ResourceDirectiveName ||
		directive.Name == consts.ResourceFieldDirectiveName ||
		directive.Name == consts.ResourceActionDirectiveName ||
		directive.Name == consts.ResourceActionFieldDirectiveName ||
		directive.Name == consts.ResourceListActionDirectiveName
}

// DirectiveToType converts directive to type
func DirectiveToType(directive *ast.DirectiveDefinition, pos *ast.Position) *ast.Definition {
	return &ast.Definition{
		Description: directive.Description,
		Kind:        ast.Object,
		Name:        directive.Name,
		Position:    pos,
		Fields:      ArgsToFields(directive.Arguments),
	}
}

// ArgsToFields converts arguments to fields
func ArgsToFields(args ast.ArgumentDefinitionList) ast.FieldList {
	fields := ast.FieldList{}
	for _, arg := range args {
		excludeDir := arg.Directives.ForName(consts.ExcludeArgumentFromType)
		if excludeDir != nil {
			excArg := excludeDir.Arguments.ForName("exclude")
			if excArg == nil {
				continue
			}
			val, err := excArg.Value.Value(nil)
			if err != nil {
				panic(err)
			}
			if val == nil || val.(bool) {
				fmt.Printf("qweqweqwe")
				continue
			}
		}
		fields = append(fields, &ast.FieldDefinition{
			Name:         arg.Name,
			Description:  arg.Description,
			Type:         arg.Type,
			Position:     arg.Position,
			DefaultValue: arg.DefaultValue,
			Directives:   arg.Directives,
		})
	}
	return fields
}

// GetDefinedObjects returns all defined objects in schema
func GetDefinedObjects(schema *ast.Schema, hasDirectives ...string) map[string]*ast.Definition {
	objs := make(map[string]*ast.Definition)
	for name, _type := range schema.Types {
		if _type.BuiltIn ||
			_type.Name == "Query" ||
			_type.Name == "Mutation" {
			continue
		}

		if len(hasDirectives) > 0 {
			var has bool
			for _, dir := range hasDirectives {
				if _type.Directives.ForName(dir) != nil {
					has = true
					break
				}
			}
			if !has {
				continue
			}
		}

		objs[name] = _type
	}
	return objs
}

// AppendFieldIfNotExists appends field to fields if it doesn't exist
func AppendFieldIfNotExists(fields []*ast.FieldDefinition, field *ast.FieldDefinition) []*ast.FieldDefinition {
	for _, f := range fields {
		if f.Name == field.Name {
			return fields
		}
	}
	return append(fields, field)
}

// IsQueryAction checks if directive is query action
func IsQueryAction(directive *ast.Directive) bool {
	resActionEnum := directive.Arguments.ForName("Action").Value.Raw
	return strings.HasSuffix(resActionEnum, "_QUERY")
}

// IsMutationAction checks if directive is mutation action
func IsMutationAction(directive *ast.Directive) bool {
	resActionEnum := directive.Arguments.ForName("Action").Value.Raw
	return strings.HasSuffix(resActionEnum, "_MUTATION")
}

// GetResourceDirectiveSingularType returns resource directive singular type
func GetResourceDirectiveSingularType(schema *ast.Schema, directive *ast.Directive) (*ast.Type, error) {
	resName := directive.Arguments.ForName("Resource").Value.Raw
	objType := FindObjectByResourceName(schema, resName)
	if objType == nil {
		return nil, fmt.Errorf("failed to find object for resource %s", resName)
	}

	return &ast.Type{
		NamedType: objType.Name,
	}, nil
}

// FindObjectByResourceName finds object by resource name
func FindObjectByResourceName(schema *ast.Schema, name string) *ast.Definition {
	objects := GetDefinedObjects(schema)
	for _, _type := range objects {
		directive := _type.Directives.ForName(consts.ResourceDirectiveName)

		if directive == nil {
			continue
		}

		resNameArg := directive.Arguments.ForName("Name")
		if resNameArg != nil && resNameArg.Value.Raw == name {
			return _type
		}

	}
	return nil
}
