package common

import (
	"fmt"
	"github.com/goxgen/goxgen/consts"
	"github.com/mitchellh/mapstructure"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

var pos = &ast.Position{Src: &ast.Source{BuiltIn: false}}

// IsXgenDirectiveDefinition checks if directive is xgen directive
func IsXgenDirectiveDefinition(directive *ast.DirectiveDefinition) bool {
	return directive.Name == consts.SchemaDefDirectiveResourceName ||
		directive.Name == consts.SchemaDefDirectiveFieldName ||
		directive.Name == consts.SchemaDefDirectiveActionName ||
		directive.Name == consts.SchemaDefDirectiveActionFieldName ||
		directive.Name == consts.SchemaDefDirectiveListActionName
}

// DirectiveToType converts directive to type
func DirectiveToType(directive *ast.DirectiveDefinition, pos *ast.Position) *ast.Definition {
	return &ast.Definition{
		Name:        directive.Name,
		Description: directive.Description,
		Kind:        ast.Object,
		Position:    pos,
		Fields:      ArgsToObjectFields(directive.Arguments),
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

func ArgsToObjectFields(args ast.ArgumentDefinitionList) ast.FieldList {
	fieldList := ArgsToFields(args)
	newFieldList, err := AnyFieldListToObjectFieldList(fieldList)
	if err != nil {
		panic(err)
	}
	return newFieldList
}

func ArgsToStruct[T any](args ast.ArgumentList, st *T) error {
	val := reflect.ValueOf(st).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		arg := args.ForName(name)
		if arg == nil {
			continue
		}
		v, err := arg.Value.Value(nil)
		if err != nil {
			return err
		}

		fd := val.Field(i)
		newValue, err := castByReflectType(v, fd.Type())
		if err != nil {
			return err
		}

		fd.Set(reflect.ValueOf(newValue))
	}
	return nil
}

func castByReflectType(a interface{}, targetType reflect.Type) (interface{}, error) {

	if a == nil && targetType == nil {
		return nil, nil
	}

	sourceType := reflect.TypeOf(a)
	sourceValue := reflect.ValueOf(a)

	// Helper function to handle pointer types
	handlePointerTypes := func(isTargetPtr bool) reflect.Type {
		if sourceType != nil && sourceType.Kind() == reflect.Ptr {
			sourceType = sourceType.Elem()
			sourceValue = sourceValue.Elem()
		}
		if isTargetPtr {
			return targetType.Elem()
		}
		return targetType
	}

	isTargetPtr := targetType.Kind() == reflect.Ptr
	targetType = handlePointerTypes(isTargetPtr)

	// Handle slices
	if sourceType != nil &&
		(sourceType.Kind() == reflect.Slice && targetType.Kind() == reflect.Slice) {
		if a == nil {
			return reflect.MakeSlice(targetType, 0, 0).Interface(), nil
		}
		slice := reflect.MakeSlice(targetType, sourceValue.Len(), sourceValue.Cap())
		for i := 0; i < sourceValue.Len(); i++ {
			elem, err := castByReflectType(sourceValue.Index(i).Interface(), targetType.Elem())
			if err != nil {
				return nil, err
			}
			if targetType.Elem().Kind() == reflect.Ptr && reflect.TypeOf(elem).Kind() != reflect.Ptr {
				temp := reflect.New(reflect.TypeOf(elem))
				temp.Elem().Set(reflect.ValueOf(elem))
				elem = temp.Interface()
			}
			slice.Index(i).Set(reflect.ValueOf(elem))
		}
		return slice.Interface(), nil
	} else if sourceType == nil && targetType.Kind() == reflect.Slice {
		return reflect.MakeSlice(targetType, 0, 0).Interface(), nil
	}

	// Handle map to struct and struct to map
	if (sourceType != nil && sourceType.Kind() == reflect.Map && targetType.Kind() == reflect.Struct) ||
		(sourceType != nil && sourceType.Kind() == reflect.Struct && targetType.Kind() == reflect.Map) {
		if a == nil {
			return reflect.MakeMap(targetType).Interface(), nil
		}
		targetValue := reflect.New(targetType).Elem()
		if err := mapstructure.Decode(a, targetValue.Addr().Interface()); err != nil {
			return nil, err
		}
		if isTargetPtr {
			temp := reflect.New(targetType)
			temp.Elem().Set(targetValue)
			return temp.Interface(), nil
		}
		return targetValue.Interface(), nil
	} else if sourceType == nil && targetType.Kind() == reflect.Map {
		return reflect.MakeMap(targetType).Interface(), nil
	}

	// Handle primitive types
	if sourceType != nil && sourceType.ConvertibleTo(targetType) {
		// Special case to disallow int to string conversions
		if sourceType.Kind() == reflect.Int && targetType.Kind() == reflect.String {
			return nil, fmt.Errorf("cannot convert %s to %s", sourceType, targetType)
		}
		converted := sourceValue.Convert(targetType).Interface()
		if isTargetPtr {
			temp := reflect.New(targetType)
			temp.Elem().Set(reflect.ValueOf(converted))
			return temp.Interface(), nil
		}
		return converted, nil
	} else if sourceType == nil {
		return nil, nil
	}

	return nil, fmt.Errorf("cannot convert %s to %s", sourceType, targetType)
}

// ToObjectDefinition converts ast.Definition with any Kind to ast.Definition with a Kind is ast.Object
func ToObjectDefinition(def ast.Definition, newName string) *ast.Definition {
	newFieldList, err := AnyFieldListToObjectFieldList(def.Fields)
	if err != nil {
		panic(err)
	}
	return &ast.Definition{
		Kind:        ast.Object,
		Name:        newName,
		Description: def.Description,
		Position:    def.Position,
		Fields:      newFieldList,
	}
}

func AnyFieldListToObjectFieldList(fields ast.FieldList) (ast.FieldList, error) {
	var newFields ast.FieldList
	for _, field := range fields {
		newField, err := AnyFieldToObjectField(*field)
		if err != nil {
			return nil, err
		}
		newFields = append(newFields, newField)
	}
	return newFields, nil
}

func AnyFieldToObjectField(field ast.FieldDefinition) (*ast.FieldDefinition, error) {
	dir := field.Directives.ForName(consts.ToObjectType)
	fieldType := field.Type
	if dir != nil {
		fieldTypeArg := dir.Arguments.ForName("type")
		if fieldTypeArg == nil {
			return nil, fmt.Errorf("failed to get type argument")
		}
		if fieldTypeArg.Value == nil {
			return nil, fmt.Errorf("failed to get type argument value")
		}
		fieldTypeValue, err := fieldTypeArg.Value.Value(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get type argument value: %w", err)
		}

		fieldTypeName, ok := fieldTypeValue.(string)
		if !ok {
			return nil, fmt.Errorf("type argument is not string")
		}
		fieldType = &ast.Type{
			NamedType: fieldTypeName,
		}
		field.Directives = slices.DeleteFunc(field.Directives, func(i *ast.Directive) bool {
			return i == dir
		})
	}
	return &ast.FieldDefinition{
		Name:         field.Name,
		Description:  field.Description,
		Type:         fieldType,
		Position:     field.Position,
		DefaultValue: field.DefaultValue,
		Directives:   field.Directives,
	}, nil
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
	resName := directive.Arguments.ForName(consts.SchemaDefActionDirectiveArgResource).Value.Raw
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
		directive := _type.Directives.ForName(consts.SchemaDefDirectiveResourceName)

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
