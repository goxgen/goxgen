# goxgen

[![GitHub license](https://img.shields.io/github/license/goxgen/goxgen)](https://github.com/goxgen/goxgen)
[![GitHub stars](https://img.shields.io/github/stars/goxgen/goxgen)](https://github.com/goxgen/goxgen/stargazers)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/goxgen/goxgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/goxgen/goxgen)](https://goreportcard.com/report/github.com/goxgen/goxgen)
[![codecov](https://codecov.io/gh/goxgen/goxgen/branch/main/graph/badge.svg?token=SDEXU6YQH9)](https://codecov.io/gh/goxgen/goxgen)

Your One-Stop Solution for GraphQL Application Generation

`goxgen` is a powerful library designed to simplify the creation of GraphQL applications.
By defining your domain and API interface through a single syntax,
You can quickly generate a fully-functional GraphQL server.
Beyond that, `goxgen` also provides support for ORM(GORM)
and a Command-Line Interface for server operations.

> Built upon the `gqlgen` framework, `goxgen` extends its
> capabilities to offer a more streamlined developer experience.

## 🌟 Features

- 📝 **Single Syntax for Domain and API:** Define your domain and API interface in GraphQL schema language.
- 📊 **GraphQL:** Schema-based application generation
- 🎛️ **ORM Support:** Seamlessly integrates with various ORM systems like GORM and ENT.
- ⚙️ **CLI Support:** Comes with a CLI tool to spin up your server application in no time.
- 📚**Domain Driven Design:** Extensible project structure
- 🛡️**Future-Ready:** Plans to roll out UI for admin back-office, along with comprehensive authentication and authorization features.

## 📦 Dependencies

- [gqlgen](https://github.com/99designs/gqlgen)
- [gorm](https://gorm.io/index.html)
- [urfave/cli](https://cli.urfave.org)

# 🚀 Quick Start

## 👣 Step-by-step guide

### 📄 Creating the necessary files

You should create two files in your project

1. Standard `gen.go` file with `go:generate` directive
    ```go
    package main
    
    //go:generate go run -mod=mod github.com/goxgen/goxgen
    
    ```
2. Xgen config file `xgenc.go`
    ```go
    package main
    
    import (
    	"context"
    	"fmt"
    	"github.com/goxgen/goxgen/plugins/cli"
    	"github.com/goxgen/goxgen/projects/gorm"
    	"github.com/goxgen/goxgen/projects/simple"
    	"github.com/goxgen/goxgen/xgen"
    )
    
    func main() {
    	xg := xgen.NewXgen(
    		xgen.WithPlugin(cli.NewPlugin()),
    		xgen.WithPackageName("github.com/goxgen/goxgen/cmd/internal/integration"),
    		xgen.WithProject(
    			"myproject",
    			simple.NewPlugin(),
    		),
    		xgen.WithProject(
    			"gormproj",
    			gorm.NewPlugin(),
    		),
    		//xgen.WithProject(
    		//	"entproj",
    		//	projects.NewEntProject(),
    		//),
    	)
    
    	err := xg.Generate(context.Background())
    	if err != nil {
    		fmt.Println(err)
    	}
    }
    
    ```
Then run `go generate` command, and goxgen will generate project structure

```shell
go generate
```

### 📁 Structure of a generated project

After running `go generate` command, goxgen will generate project structure like this

```shell
demoproj/
|-- entproj/
|-- gormproj/
|   |-- generate/
|   |-- generated_xgen_directives.graphql
|   |-- generated_xgen_introspection.go
|   |-- generated_xgen_introspection.graphql
|   |-- generated_xgen_project_handlers.go
|   |-- graphql.config.yml
|   |-- resolver.go
|   |-- schema.main.go
|   |-- schema.resolver.go
|-- myproject/
|-- .gitignore
|-- gen.go
|-- generated_xgen_cli.go
|-- go.mod
|-- xgenc.go
```

### 📑 Providing schema

You should provide a schema for each project and run `go generate` again.

All schema files in xgen has this format `schema.{some_name}.graphql`, for example `schema.main.graphql`


Let's focus on `gormproj`, which uses the GORM ORM.
The connection to the GORM database can be configured from the gqlgen standard `resolver.go` file in the `gormproj` directory.

> `resolver.go` is designed to support your custom dependency injection (DI) and any services you've provided.
```go
package gormproj

import (
	"fmt"
	"github.com/goxgen/goxgen/cmd/internal/integration/gormproj/generated"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}

	// Open the database connection
	db, err := gorm.Open(sqlite.Open("./cmd/internal/integration/gormproj.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&generated.Car{},
		&generated.User{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	r.DB = db
	return r, nil
}

```

### Example of schema file `schema.main.graphql`
```graphql
type User
@Resource(Name: "user", Primary: true, Route: "user", DB: {Table: "user"})
{
    id: ID! @Field(Label: "ID", Description: "ID of the todo", DB: {Column: "id", PrimaryKey: true})
    name: String! @Field(Label: "Text", Description: "Text of the todo", DB: {Column: "name", Unique: true})
    cars: [Car!]! @Field(Label: "Cars", Description: "Cars of the todo", DB: {})
}

type Car
@Resource(Name: "car", Primary: true, Route: "car", DB: {Table: "car"})
{
    id: ID! @Field(Label: "ID", Description: "ID of the todo", DB: {Column: "id", PrimaryKey: true})
    make: String! @Field(Label: "Make", Description: "Car make", DB: {Column: "make"})
    done: Boolean! @Field(Label: "Done", Description: "Done of the todo", DB: {Column: "done"})
    user: User! @Field(Label: "User", Description: "User of the todo", DB: {})
}

input CarInput
@Action(Resource: "car", Action: CREATE_MUTATION, Route: "new", SchemaFieldName: "new_car")
@Action(Resource: "car", Action: UPDATE_MUTATION, Route: "update", SchemaFieldName: "update_car")
{
    id: ID @ActionField(Label: "ID", Description: "ID of the todo")
    make: String @ActionField(Label: "Make", Description: "Text of the todo")
    done: Boolean @ActionField(Label: "Done", Description: "Done of the todo")
    user: ID @ActionField(Label: "User", Description: "User of the todo")
}

input NewUser
@Action(Resource: "user", Action: CREATE_MUTATION, Route: "new", SchemaFieldName: "new_user")
{
    name: String! @ActionField(Label: "Name", Description: "Name")
    cars: [CarInput!] @ActionField(Label: "Cars", Description: "Cars of the todo")
}

input DeleteUsers
@ListAction(Resource: "user", Action: BATCH_DELETE_MUTATION, Route: "delete", SchemaFieldName: "delete_users")
{
    ids: [ID!] @ActionField(Label: "IDs", Description: "IDs of users")
}

input ListUser
@ListAction(Resource: "user", Action: BROWSE_QUERY, Route: "list", SchemaFieldName: "list_user")
{
    id: ID @ActionField(Label: "ID", Description: "ID")
    name: String @ActionField(Label: "Name", Description: "Name")
}

input ListCars
@ListAction(Resource: "car", Action: BROWSE_QUERY, Route: "list", SchemaFieldName: "list_cars")
{
    id: ID @ActionField(Label: "ID", Description: "ID")
    userId: ID @ActionField(Label: "User ID", Description: "User ID")
    make: String @ActionField(Label: "Make", Description: "Make")
}
```

The directives used in the example above are standard `xgen` directives, intended to provide metadata.

* `Resource` - Entity or object or thing
* `Field` - Field of resource
* `Action` - Action that can be done for single resource
* `ListAction` - Action that can be done for bulk resources
* `ActionField` - Field of action or list action

The definitions of these directives are located in the `generated_xgen_directives.graphql` file.
```graphql
"""This directive is used to mark the object as a resource"""
directive @Resource(Name: String!, Route: String, Primary: Boolean, DB: XgenResourceDbConfigInput @ExcludeArgumentFromType) on OBJECT
"""This directive is used to mark the object as a resource action"""
directive @Action(Resource: String!, Action: XgenResourceActionType!, Route: String, SchemaFieldName: String) repeatable on INPUT_OBJECT
"""This directive is used to mark the object as a resource list action"""
directive @ListAction(Resource: String!, Action: XgenResourceListActionType!, Route: String, Pagination: Boolean, SchemaFieldName: String) repeatable on INPUT_OBJECT
"""This directive is used to exclude the argument from the type"""
directive @ExcludeArgumentFromType(exclude: Boolean) on ARGUMENT_DEFINITION
"""This directive is used to mark the object as a resource field"""
directive @Field(Label: String, Description: String, DB: XgenResourceFieldDbConfigInput @ExcludeArgumentFromType) on FIELD_DEFINITION
"""This directive is used to mark the object as a resource field"""
directive @ActionField(Label: String, Description: String) on INPUT_FIELD_DEFINITION
enum XgenResourceActionType {
  CREATE_MUTATION
  READ_QUERY
  UPDATE_MUTATION
  DELETE_MUTATION
}
enum XgenResourceListActionType {
  BROWSE_QUERY
  BATCH_DELETE_MUTATION
}
input XgenPaginationInput {
  page: Int!
  limit: Int!
}
input XgenCursorPaginationInput {
  first: Int!
  after: String
  last: Int!
  before: String
}
input XgenResourceDbConfigInput {
  Table: String
}
input XgenResourceFieldDbConfigInput {
  Column: String
  PrimaryKey: Boolean
  AutoIncrement: Boolean
  Unique: Boolean
  NotNull: Boolean
  Index: Boolean
  UniqueIndex: Boolean
  Size: Int
  Precision: Int
  Type: String
  Scale: Int
  AutoIncrementIncrement: Int
}
```

After writing a custom schema You should run again `gogen` command.

```shell
go generate
```

After regenerating the code, the `schema.resolver.go` file will be updated based on your schema. 
You can expect to see changes similar to the following:

```go
package gormproj

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/goxgen/goxgen/cmd/internal/integration/gormproj/generated"
	"github.com/goxgen/goxgen/utils"
	"gorm.io/gorm/clause"
)

// DeleteUsers is the resolver for the delete_users field.
func (r *mutationResolver) DeleteUsers(ctx context.Context, input *generated.DeleteUsers) ([]*generated.User, error) {
	var users []*generated.User
	r.DB.Delete(&users, input.Ids)
	return users, nil
}

// NewUser is the resolver for the new_user field.
func (r *mutationResolver) NewUser(ctx context.Context, input *generated.NewUser) (*generated.User, error) {
	cars := make([]*generated.Car, len(input.Cars))
	for i, car := range input.Cars {
		cars[i] = &generated.Car{
			ID:   utils.Deref(car.ID),
			Make: utils.Deref(car.Make),
			Done: utils.Deref(car.Done),
		}
	}
	user := &generated.User{
		Name: input.Name,
		Cars: cars,
	}
	res := r.DB.Preload(clause.Associations).Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

// NewCar is the resolver for the new_car field.
func (r *mutationResolver) NewCar(ctx context.Context, input *generated.CarInput) (*generated.Car, error) {
	panic(fmt.Errorf("not implemented: NewCar - new_car"))
}

// UpdateCar is the resolver for the update_car field.
func (r *mutationResolver) UpdateCar(ctx context.Context, input *generated.CarInput) (*generated.Car, error) {
	panic(fmt.Errorf("not implemented: UpdateCar - update_car"))
}

// XgenIntrospection is the resolver for the _xgen_introspection field.
func (r *queryResolver) XgenIntrospection(ctx context.Context) (*generated.XgenIntrospection, error) {
	return r.Resolver.XgenIntrospection()
}

// ListCars is the resolver for the list_cars field.
func (r *queryResolver) ListCars(ctx context.Context, input *generated.ListCars) ([]*generated.Car, error) {
	var cars []*generated.Car
	res := r.DB.Where(&[]*generated.Car{
		{
			ID:     utils.Deref(input.ID),
			UserID: utils.Deref(input.UserID),
		},
	}).Find(&cars)

	if res.Error != nil {
		return nil, res.Error
	}
	return cars, nil
}

// ListUser is the resolver for the list_user field.
func (r *queryResolver) ListUser(ctx context.Context, input *generated.ListUser) ([]*generated.User, error) {
	var users []*generated.User
	res := r.DB.Where(&[]*generated.User{
		{
			ID:   utils.Deref(input.ID),
			Name: utils.Deref(input.Name),
		},
	}).Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

```

You can add your own implementation for each function in the updated `schema.resolver.go` file.
For more information,
You can read the [gqlgen documentation](https://gqlgen.com/getting-started/#implement-the-resolvers). 


You can see that some functions in example are implemented, some are not. 
So take a look at the implemented functions, for example `NewUser`, `ListUser` or `DeleteUsers`.
In those functions, you can see that the `r.DB` instance is used, 
which is provided from the `resolver.go` file.

```go
package gormproj

import (
	"fmt"
	"github.com/goxgen/goxgen/cmd/internal/integration/gormproj/generated"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}

	// Open the database connection
	db, err := gorm.Open(sqlite.Open("./cmd/internal/integration/gormproj.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&generated.Car{},
		&generated.User{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	r.DB = db
	return r, nil
}

```

Great, you're all set to launch your GraphQL application.

To start the server using the xgen CLI plugin, you can run the following command:

```shell
go run generated_xgen_cli.go serve --gql-playground-enabled gormproj
```

This will initialize and start your GraphQL server, making it ready to handle incoming requests.

The output from the xgen CLI will provide information about the server endpoints. Additionally, logs will be written to this output during the server's runtime, giving you insights into its operation.

```shell
2023-08-30T13:34:58.750+0400    INFO    gormproj/generated_xgen_project_handlers.go:102 Serving graphql playground      {"url": "http://localhost:80/playground"}
2023-08-30T13:34:58.750+0400    INFO    gormproj/generated_xgen_project_handlers.go:113 Serving graphql                 {"url": "http://localhost:80/query"}
```

> For more information about the xgen CLI, you can run the following command: 
> 
> `go run generated_xgen_cli.go help`
> 
> This will display a list of available commands, options, and descriptions to help you navigate the xgen CLI more effectively.

You can copy the URL `http://localhost:80/playground` from the logs 
and open it in your browser to access the GraphQL playground. 
This interface will allow you to test queries, mutations, and subscriptions in real-time.

Then we see graphql playground, let's run some mutation query to add two new users

```graphql
mutation{
    user1: new_user(input: {name: "My user 1", cars: {make:"BMW"}}){
        id
        name
        cars {
            make
            id
        }
    }
    user2: new_user(input: {name: "My user 2", cars: {make:"BMW"}}){
        id
        name
        cars {
            make
            id
        }
    }
}
```

After execution of this mutation, graphql should be return result like this

```json
{
  "data": {
    "user1": {
      "id": 1,
      "name": "My user 1",
      "cars": [
        {
          "make": "BMW",
          "id": 12
        }
      ]
    },
    "user2": {
      "id": 2,
      "name": "My user 2",
      "cars": [
        {
          "make": "BMW",
          "id": 13
        }
      ]
    }
  }
}
```

One more example, let's list our new users by query
```graphql
query{
    list_user(input: {}){
        id
        name
    }
}
```

The result of this query should be like this
