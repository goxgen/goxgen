// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"fmt"
	"io"
	"strconv"
)

// This directive is used to mark the object as a resource action
type Action struct {
	Resource        string                 `json:"Resource"`
	Action          XgenResourceActionType `json:"Action"`
	Route           *string                `json:"Route,omitempty"`
	SchemaFieldName *string                `json:"SchemaFieldName,omitempty"`
}

type ActionAnnotationSingle struct {
	Name  *string `json:"name,omitempty"`
	Value *Action `json:"value,omitempty"`
}

// This directive is used to mark the object as a resource field
type ActionField struct {
	Label       *string `json:"Label,omitempty"`
	Description *string `json:"Description,omitempty"`
}

// This directive is used to mark the object as a resource field
type Field struct {
	Label       *string `json:"Label,omitempty"`
	Description *string `json:"Description,omitempty"`
}

// This directive is used to mark the object as a resource list action
type ListAction struct {
	Resource        string                     `json:"Resource"`
	Action          XgenResourceListActionType `json:"Action"`
	Route           *string                    `json:"Route,omitempty"`
	Pagination      *bool                      `json:"Pagination,omitempty"`
	SchemaFieldName *string                    `json:"SchemaFieldName,omitempty"`
}

type ListActionAnnotationSingle struct {
	Name  *string     `json:"name,omitempty"`
	Value *ListAction `json:"value,omitempty"`
}

// This directive is used to mark the object as a resource
type Resource struct {
	Name    string  `json:"Name"`
	Route   *string `json:"Route,omitempty"`
	Primary *bool   `json:"Primary,omitempty"`
}

type ResourceAnnotationSingle struct {
	Name  *string   `json:"name,omitempty"`
	Value *Resource `json:"value,omitempty"`
}

type XgenAnnotationMap struct {
	Resource   []*ResourceAnnotationSingle   `json:"Resource"`
	ListAction []*ListActionAnnotationSingle `json:"ListAction"`
	Action     []*ActionAnnotationSingle     `json:"Action"`
}

type XgenCursorPaginationInput struct {
	First  int     `json:"first"`
	After  *string `json:"after,omitempty"`
	Last   int     `json:"last"`
	Before *string `json:"before,omitempty"`
}

type XgenCursorPaginationInputXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenFieldDef struct {
	Field       *Field       `json:"Field,omitempty"`
	ActionField *ActionField `json:"ActionField,omitempty"`
}

type XgenIntrospection struct {
	Annotation *XgenAnnotationMap `json:"annotation,omitempty"`
	Object     *XgenObjectMap     `json:"object,omitempty"`
}

type XgenObjectDefinition struct {
	Resource   *Resource   `json:"Resource,omitempty"`
	ListAction *ListAction `json:"ListAction,omitempty"`
	Action     *Action     `json:"Action,omitempty"`
}

type XgenObjectField struct {
	Name       *string       `json:"name,omitempty"`
	Definition *XgenFieldDef `json:"definition,omitempty"`
}

type XgenObjectMap struct {
	XgenCursorPaginationInput      *XgenCursorPaginationInputXgenDef      `json:"XgenCursorPaginationInput,omitempty"`
	XgenResourceActionType         *XgenResourceActionTypeXgenDef         `json:"XgenResourceActionType,omitempty"`
	XgenResourceFieldDbConfigInput *XgenResourceFieldDbConfigInputXgenDef `json:"XgenResourceFieldDbConfigInput,omitempty"`
	XgenPaginationInput            *XgenPaginationInputXgenDef            `json:"XgenPaginationInput,omitempty"`
	XgenResourceDbConfigInput      *XgenResourceDbConfigInputXgenDef      `json:"XgenResourceDbConfigInput,omitempty"`
	XgenResourceListActionType     *XgenResourceListActionTypeXgenDef     `json:"XgenResourceListActionType,omitempty"`
}

type XgenPaginationInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type XgenPaginationInputXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenResourceActionTypeXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenResourceDbConfigInput struct {
	Table *string `json:"Table,omitempty"`
}

type XgenResourceDbConfigInputXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenResourceFieldDbConfigInput struct {
	Column                 *string `json:"Column,omitempty"`
	PrimaryKey             *bool   `json:"PrimaryKey,omitempty"`
	AutoIncrement          *bool   `json:"AutoIncrement,omitempty"`
	Unique                 *bool   `json:"Unique,omitempty"`
	NotNull                *bool   `json:"NotNull,omitempty"`
	Index                  *bool   `json:"Index,omitempty"`
	UniqueIndex            *bool   `json:"UniqueIndex,omitempty"`
	Size                   *int    `json:"Size,omitempty"`
	Precision              *int    `json:"Precision,omitempty"`
	Type                   *string `json:"Type,omitempty"`
	Scale                  *int    `json:"Scale,omitempty"`
	AutoIncrementIncrement *int    `json:"AutoIncrementIncrement,omitempty"`
}

type XgenResourceFieldDbConfigInputXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenResourceListActionTypeXgenDef struct {
	Object *XgenObjectDefinition `json:"object,omitempty"`
	Field  []*XgenObjectField    `json:"field"`
}

type XgenResourceActionType string

const (
	XgenResourceActionTypeCreateMutation XgenResourceActionType = "CREATE_MUTATION"
	XgenResourceActionTypeReadQuery      XgenResourceActionType = "READ_QUERY"
	XgenResourceActionTypeUpdateMutation XgenResourceActionType = "UPDATE_MUTATION"
	XgenResourceActionTypeDeleteMutation XgenResourceActionType = "DELETE_MUTATION"
)

var AllXgenResourceActionType = []XgenResourceActionType{
	XgenResourceActionTypeCreateMutation,
	XgenResourceActionTypeReadQuery,
	XgenResourceActionTypeUpdateMutation,
	XgenResourceActionTypeDeleteMutation,
}

func (e XgenResourceActionType) IsValid() bool {
	switch e {
	case XgenResourceActionTypeCreateMutation, XgenResourceActionTypeReadQuery, XgenResourceActionTypeUpdateMutation, XgenResourceActionTypeDeleteMutation:
		return true
	}
	return false
}

func (e XgenResourceActionType) String() string {
	return string(e)
}

func (e *XgenResourceActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = XgenResourceActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid XgenResourceActionType", str)
	}
	return nil
}

func (e XgenResourceActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type XgenResourceListActionType string

const (
	XgenResourceListActionTypeBrowseQuery         XgenResourceListActionType = "BROWSE_QUERY"
	XgenResourceListActionTypeBatchDeleteMutation XgenResourceListActionType = "BATCH_DELETE_MUTATION"
)

var AllXgenResourceListActionType = []XgenResourceListActionType{
	XgenResourceListActionTypeBrowseQuery,
	XgenResourceListActionTypeBatchDeleteMutation,
}

func (e XgenResourceListActionType) IsValid() bool {
	switch e {
	case XgenResourceListActionTypeBrowseQuery, XgenResourceListActionTypeBatchDeleteMutation:
		return true
	}
	return false
}

func (e XgenResourceListActionType) String() string {
	return string(e)
}

func (e *XgenResourceListActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = XgenResourceListActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid XgenResourceListActionType", str)
	}
	return nil
}

func (e XgenResourceListActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
