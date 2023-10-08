package mapper

import (
	"go.uber.org/zap"
	"reflect"
	"strings"
)

// Mapper is a struct mapper.
type Mapper struct {
	logger      *zap.Logger
	maptoTagKey string
}

// Option is a mapper option.
type Option func(*Mapper) error

// New creates a new mapper.
func New(opts ...Option) (*Mapper, error) {
	m := &Mapper{
		logger:      zap.NewNop(),
		maptoTagKey: "mapto",
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// WithLogger sets the logger for the mapper.
func WithLogger(logger *zap.Logger) Option {
	return func(m *Mapper) error {
		m.logger = logger
		return nil
	}
}

// Map maps fields from a source struct to a destination struct based on "mapto" tags.
func (m *Mapper) Map(src interface{}, dest interface{}) error {
	m.logger.Info("Mapping",
		zap.String("src", reflect.TypeOf(src).String()),
		zap.String("dest", reflect.TypeOf(dest).String()),
	)

	srcRef := reflect.ValueOf(src)

	if srcRef.Kind() == reflect.Ptr && srcRef.IsNil() {
		return nil // Skip if the source is nil
	}

	destRef := reflect.ValueOf(dest)

	if destRef.Kind() != reflect.Ptr {
		return NewDestinationNotPointerError("destination must be a pointer")
	}

	srcValue, destValue := m.dereferencePointers(srcRef, destRef)

	if srcValue.Kind() != reflect.Struct {
		return NewSourceNotStructError("source must be a struct")
	}

	if destValue.Kind() != reflect.Struct {
		return NewDestinationNotStructError("destination must be a struct")
	}

	return m.mapFields(srcValue, destValue)
}

// dereferencePointers dereferences pointers and returns the underlying values.
func (m *Mapper) dereferencePointers(srcValue, destValue reflect.Value) (reflect.Value, reflect.Value) {
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	if destValue.Kind() == reflect.Ptr {
		if destValue.IsNil() {
			destValue.Set(reflect.New(destValue.Type().Elem()))
		}
		destValue = destValue.Elem()
	}
	return srcValue, destValue
}

// mapFields maps fields from srcValue to destValue based on "mapto" tags.
func (m *Mapper) mapFields(srcValue, destValue reflect.Value) error {
	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		maptoTag := field.Tag.Get(m.maptoTagKey)
		if maptoTag == "" {
			continue
		}

		maptoFields := strings.Split(maptoTag, ",")
		for _, maptoField := range maptoFields {
			destStructName, destFieldName, err := m.parseTag(maptoField, field, srcValue)
			if err != nil {
				return err
			}

			if destStructName != destValue.Type().Name() {
				continue
			}

			srcField := srcValue.Field(i)
			// Skip if the source is nil
			if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
				continue // Skip if the source is nil
			}

			destField := destValue.FieldByName(destFieldName)

			if !destField.IsValid() {
				m.logger.Error("Destination field not found",
					zap.String("field", destFieldName),
				)
				return NewDestinationFieldNotFoundError("destination field %s not found", destFieldName)
			}

			// Initialize nested pointer fields in the destination
			if destField.Kind() == reflect.Ptr && destField.IsNil() {
				destField.Set(reflect.New(destField.Type().Elem()))
			}

			// Handle slice fields
			if srcField.Kind() == reflect.Slice {
				if err := m.handleSliceFields(srcField, destField); err != nil {
					return err
				}
				continue
			}

			m.logger.Info("Mapping field",
				zap.String("field", field.Name),
				zap.Any("srcField", srcField),
				zap.Any("destField", destField),
			)

			// Type mismatch check
			if srcField.IsValid() && destField.IsValid() {
				if !destField.CanSet() {
					return NewDestinationFieldNotSettableError("destination field %s(%s) is not settable", destFieldName, destField.Type())
				}

				if srcField.Type() == destField.Type() {
					destField.Set(srcField)
				} else if srcField.Kind() == reflect.Ptr && srcField.Elem().Type() == destField.Type() {
					// Special case: source is a pointer, and destination is a value of the same base type
					//dereferencePointers(srcField, destField)
					if !srcField.IsNil() {
						destField.Set(srcField.Elem())
					}
				} else if destField.Kind() == reflect.Ptr && destField.Elem().Type() == srcField.Type() {
					// Special case: destination is a pointer, and source is a value of the same base type
					//dereferencePointers(srcField, destField)
					if destField.IsNil() {
						destField.Set(reflect.New(destField.Type().Elem()))
					}
					destField.Elem().Set(srcField)
				} else if srcField.Kind() == reflect.Ptr && destField.Kind() == reflect.Ptr &&
					srcField.Elem().Kind() == reflect.Struct && destField.Elem().Kind() == reflect.Struct {
					// Special case: both fields are pointers to different struct types
					if srcField.IsNil() {
						continue // Skip if the source is nil
					}
					if destField.IsNil() {
						destField.Set(reflect.New(destField.Type().Elem()))
					}

					if destField.Kind() != reflect.Ptr {
						destField = destField.Addr()
					}

					err := m.Map(srcField.Elem().Interface(), destField.Interface())
					if err != nil {
						return err
					}
				} else {
					return NewFieldsTypesMismatchError("type mismatch for field %s: src type %s, dest type %s", field.Name, srcField.Type(), destField.Type())
				}
			} else {
				if !srcField.IsValid() {
					err = NewSourceFieldIsNotValidError("source field %s is not valid", field.Name)
				}
				if !destField.IsValid() {
					err = NewDestFieldIsNotValidError("destination field %s is not valid", field.Name)
				}
				return err
			}
		}
	}
	return nil
}

// parseTag parses the "mapto" tag and returns the destination struct and field names.
func (m *Mapper) parseTag(maptoField string, field reflect.StructField, srcValue reflect.Value) (string, string, error) {
	maptoFieldParts := strings.Split(maptoField, ".")
	if len(maptoFieldParts) != 2 {
		return "", "", NewInvalidMaptoTagError("invalid mapto tag for field %s in struct %s", field.Name, srcValue.Type().Name())
	}
	return maptoFieldParts[0], maptoFieldParts[1], nil
}

// handleSliceFields handles the mapping of slice fields.
func (m *Mapper) handleSliceFields(srcField, destField reflect.Value) error {
	destField.Set(reflect.MakeSlice(destField.Type(), srcField.Len(), srcField.Cap()))
	for j := 0; j < srcField.Len(); j++ {
		srcElem := srcField.Index(j)
		destElem := reflect.New(destField.Type().Elem().Elem()) // Create a new element for the destination slice

		if srcElem.Kind() == reflect.Ptr && srcElem.IsNil() {
			destField.Index(j).Set(destElem) // Set the newly created (but uninitialized) element
			continue
		}

		if srcElem.Kind() == reflect.Ptr {
			srcElem = srcElem.Elem()
		}

		err := m.Map(srcElem.Interface(), destElem.Interface())
		if err != nil {
			return err
		}

		destField.Index(j).Set(destElem) // Set the newly created and initialized element
	}
	return nil
}
