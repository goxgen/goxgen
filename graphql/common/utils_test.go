package common

import (
	"github.com/goxgen/goxgen/utils"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SomeStruct struct {
	Field1 string
	Field2 int
}

func TestCastByReflectType(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		targetType  reflect.Type
		expected    interface{}
		expectError bool
	}{
		{
			name:       "int to float64",
			input:      42,
			targetType: reflect.TypeOf(float64(0)),
			expected:   float64(42),
		},
		{
			name:       "[]*string to []string",
			input:      []*string{new(string), new(string)},
			targetType: reflect.TypeOf([]string{}),
			expected:   []string{"", ""},
		},
		{
			name:       "[]string to []*string",
			input:      []string{"", ""},
			targetType: reflect.TypeOf([]*string{}),
			expected:   []*string{new(string), new(string)},
		},
		{
			name:       "[]interface{} to []*int",
			input:      []interface{}{1, 2, 3},
			targetType: reflect.TypeOf([]*int{}),
			expected:   []*int{utils.IntP(1), utils.IntP(2), utils.IntP(3)},
		},
		{
			name:       "map[string]interface{} to SomeStruct",
			input:      map[string]interface{}{"Field1": "hello", "Field2": 42},
			targetType: reflect.TypeOf(SomeStruct{}),
			expected:   SomeStruct{"hello", 42},
		},
		{
			name:       "*map[string]interface{} to *SomeStruct",
			input:      &map[string]interface{}{"Field1": "hello", "Field2": 42},
			targetType: reflect.TypeOf(&SomeStruct{}),
			expected:   &SomeStruct{"hello", 42},
		},
		{
			name:        "incompatible types",
			input:       42,
			targetType:  reflect.TypeOf("string"),
			expectError: true,
		},
		{
			name:       "float64 to int",
			input:      42.0,
			targetType: reflect.TypeOf(int(0)),
			expected:   int(42),
		},
		{
			name:        "bool to string disallowed",
			input:       true,
			targetType:  reflect.TypeOf(""),
			expectError: true,
		},
		{
			name:       "struct to map",
			input:      SomeStruct{"hello", 42},
			targetType: reflect.TypeOf(map[string]interface{}{}),
			expected:   map[string]interface{}{"Field1": "hello", "Field2": 42},
		},
		{
			name:       "nil slice to empty slice",
			input:      nil,
			targetType: reflect.TypeOf([]int{}),
			expected:   []int{},
		},
		{
			name:       "nil map to empty map",
			input:      nil,
			targetType: reflect.TypeOf(map[string]int{}),
			expected:   map[string]int{},
		},
		{
			name:       "nil pointer to nil",
			input:      nil,
			targetType: reflect.TypeOf(new(int)),
			expected:   nil,
		},
		{
			name:        "empty string to nil",
			input:       "",
			targetType:  reflect.TypeOf(new(int)),
			expectError: true,
		},
		{
			name:        "nil to nil",
			input:       nil,
			targetType:  reflect.TypeOf(nil),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "nil to int",
			input:       nil,
			targetType:  reflect.TypeOf(int(0)),
			expected:    nil,
			expectError: false,
		},
		{
			name:        "nil slice to nil slice",
			input:       []int(nil),
			targetType:  reflect.TypeOf([]int(nil)),
			expected:    []int(nil),
			expectError: false,
		},
		{
			name:        "nil map to nil map",
			input:       map[string]int(nil),
			targetType:  reflect.TypeOf(map[string]int(nil)),
			expected:    map[string]int(nil),
			expectError: false,
		},
		{
			name:        "int to string disallowed",
			input:       42,
			targetType:  reflect.TypeOf(""),
			expectError: true,
		},
		{
			name:        "float64 to int",
			input:       float64(42.0),
			targetType:  reflect.TypeOf(int(0)),
			expected:    int(42),
			expectError: false,
		},
		{
			name:        "struct to map",
			input:       SomeStruct{"hello", 42},
			targetType:  reflect.TypeOf(map[string]interface{}{}),
			expected:    map[string]interface{}{"Field1": "hello", "Field2": 42},
			expectError: false,
		},
		{
			name:        "map to struct with extra fields",
			input:       map[string]interface{}{"Field1": "hello", "Field2": 42, "Field3": "extra"},
			targetType:  reflect.TypeOf(SomeStruct{}),
			expected:    SomeStruct{"hello", 42},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := castByReflectType(tt.input, tt.targetType)
			if tt.expectError {
				t.Logf("result: %s, err: %s", result, err)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.targetType == nil {
					assert.Nil(t, result)
					return
				} else if tt.targetType.Kind() == reflect.Slice {
					// Special handling for slice types
					assert.ElementsMatch(t, tt.expected, result)
				} else {
					assert.Equal(t, tt.expected, result)
				}
			}
		})
	}
}
