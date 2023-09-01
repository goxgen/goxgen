// Code generated by github.com/goxgen/goxgen, DO NOT EDIT.

package gormproj

import (
    _ "embed"
    "encoding/json"
    "github.com/goxgen/goxgen/internal/integration/gormproj/generated"
)

// XgenIntrospection is the resolver for the XgenIntrospection field.
func (r *Resolver) XgenIntrospection() (*generated.XgenIntrospection, error) {
    var intr generated.XgenIntrospection
    err := json.Unmarshal(
		[]byte(`{"annotation":{"Action":[{"name":"CarInput","value":{"Action":"CREATE_MUTATION","Resource":"car","Route":"new","SchemaFieldName":"new_car"}},{"name":"CarInput","value":{"Action":"UPDATE_MUTATION","Resource":"car","Route":"update","SchemaFieldName":"update_car"}},{"name":"NewUser","value":{"Action":"CREATE_MUTATION","Resource":"user","Route":"new","SchemaFieldName":"new_user"}}],"ListAction":[{"name":"ListCars","value":{"Action":"BROWSE_QUERY","Resource":"car","Route":"list","SchemaFieldName":"list_cars"}},{"name":"DeleteUsers","value":{"Action":"BATCH_DELETE_MUTATION","Resource":"user","Route":"delete","SchemaFieldName":"delete_users"}},{"name":"ListUser","value":{"Action":"BROWSE_QUERY","Resource":"user","Route":"list","SchemaFieldName":"list_user"}}],"Resource":[{"name":"Car","value":{"DB":{"Table":"car"},"Name":"car","Primary":true,"Route":"car"}},{"name":"User","value":{"DB":{"Table":"user"},"Name":"user","Primary":true,"Route":"user"}}]},"object":{"Car":{"field":[{"definition":{"Field":{"DB":{"Column":"id","PrimaryKey":true},"Description":"ID of the todo","Label":"ID"}},"name":"id"},{"definition":{"Field":{"DB":{"Column":"make"},"Description":"Car make","Label":"Make"}},"name":"make"},{"definition":{"Field":{"DB":{"Column":"done"},"Description":"Done of the todo","Label":"Done"}},"name":"done"},{"definition":{"Field":{"DB":{},"Description":"User of the todo","Label":"User"}},"name":"user"}],"object":{"Resource":{"DB":{"Table":"car"},"Name":"car","Primary":true,"Route":"car"}}},"CarInput":{"field":[{"definition":{"ActionField":{"Description":"ID of the todo","Label":"ID"}},"name":"id"},{"definition":{"ActionField":{"Description":"Text of the todo","Label":"Make"}},"name":"make"},{"definition":{"ActionField":{"Description":"Done of the todo","Label":"Done"}},"name":"done"},{"definition":{"ActionField":{"Description":"User of the todo","Label":"User"}},"name":"user"}],"object":{"Action":{"Action":"UPDATE_MUTATION","Resource":"car","Route":"update","SchemaFieldName":"update_car"}}},"DeleteUsers":{"field":[{"definition":{"ActionField":{"Description":"IDs of users","Label":"IDs"}},"name":"ids"}],"object":{"ListAction":{"Action":"BATCH_DELETE_MUTATION","Resource":"user","Route":"delete","SchemaFieldName":"delete_users"}}},"ListCars":{"field":[{"definition":{"ActionField":{"Description":"ID","Label":"ID"}},"name":"id"},{"definition":{"ActionField":{"Description":"User ID","Label":"User ID"}},"name":"userId"},{"definition":{"ActionField":{"Description":"Make","Label":"Make"}},"name":"make"}],"object":{"ListAction":{"Action":"BROWSE_QUERY","Resource":"car","Route":"list","SchemaFieldName":"list_cars"}}},"ListUser":{"field":[{"definition":{"ActionField":{"Description":"ID","Label":"ID"}},"name":"id"},{"definition":{"ActionField":{"Description":"Name","Label":"Name"}},"name":"name"}],"object":{"ListAction":{"Action":"BROWSE_QUERY","Resource":"user","Route":"list","SchemaFieldName":"list_user"}}},"NewUser":{"field":[{"definition":{"ActionField":{"Description":"Name","Label":"Name"}},"name":"name"},{"definition":{"ActionField":{"Description":"Cars of the todo","Label":"Cars"}},"name":"cars"}],"object":{"Action":{"Action":"CREATE_MUTATION","Resource":"user","Route":"new","SchemaFieldName":"new_user"}}},"User":{"field":[{"definition":{"Field":{"DB":{"Column":"id","PrimaryKey":true},"Description":"ID of the todo","Label":"ID"}},"name":"id"},{"definition":{"Field":{"DB":{"Column":"name","Unique":true},"Description":"Text of the todo","Label":"Text"}},"name":"name"},{"definition":{"Field":{"DB":{},"Description":"Cars of the todo","Label":"Cars"}},"name":"cars"}],"object":{"Resource":{"DB":{"Table":"user"},"Name":"user","Primary":true,"Route":"user"}}},"XgenCursorPaginationInput":{"field":[{"definition":{},"name":"first"},{"definition":{},"name":"after"},{"definition":{},"name":"last"},{"definition":{},"name":"before"}],"object":{}},"XgenPaginationInput":{"field":[{"definition":{},"name":"page"},{"definition":{},"name":"limit"}],"object":{}},"XgenResourceActionType":{"field":[],"object":{}},"XgenResourceDbConfigInput":{"field":[{"definition":{},"name":"Table"}],"object":{}},"XgenResourceFieldDbConfigInput":{"field":[{"definition":{},"name":"Column"},{"definition":{},"name":"PrimaryKey"},{"definition":{},"name":"AutoIncrement"},{"definition":{},"name":"Unique"},{"definition":{},"name":"NotNull"},{"definition":{},"name":"Index"},{"definition":{},"name":"UniqueIndex"},{"definition":{},"name":"Size"},{"definition":{},"name":"Precision"},{"definition":{},"name":"Type"},{"definition":{},"name":"Scale"},{"definition":{},"name":"AutoIncrementIncrement"}],"object":{}},"XgenResourceListActionType":{"field":[],"object":{}}},"resource":{"car":{"actions":[{"Action":"CREATE_MUTATION","Resource":"car","Route":"new","SchemaFieldName":"new_car"},{"Action":"UPDATE_MUTATION","Resource":"car","Route":"update","SchemaFieldName":"update_car"},{"Action":"BROWSE_QUERY","Resource":"car","Route":"list","SchemaFieldName":"list_cars"}],"objectName":"Car","properties":{"DB":{"Table":"car"},"Name":"car","Primary":true,"Route":"car"}},"user":{"actions":[{"Action":"BROWSE_QUERY","Resource":"user","Route":"list","SchemaFieldName":"list_user"},{"Action":"BATCH_DELETE_MUTATION","Resource":"user","Route":"delete","SchemaFieldName":"delete_users"},{"Action":"CREATE_MUTATION","Resource":"user","Route":"new","SchemaFieldName":"new_user"}],"objectName":"User","properties":{"DB":{"Table":"user"},"Name":"user","Primary":true,"Route":"user"}}}}`),
		&intr,
	)
    if err != nil {
        panic(err)
        return nil, err
    }
    return &intr, nil
}