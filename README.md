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
- 📚 **Domain Driven Design:** Extensible project structure
- 🛡️ **Future-Ready:** Plans to roll out UI for admin back-office, along with comprehensive authentication and authorization features.

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
    	"github.com/goxgen/goxgen/projects/basic"
    	"github.com/goxgen/goxgen/projects/gorm"
    	"github.com/goxgen/goxgen/xgen"
    )
    
    func main() {
    	xg := xgen.NewXgen(
    		xgen.WithPlugin(cli.NewPlugin()),
    		xgen.WithPackageName("github.com/goxgen/goxgen/cmd/internal/integration"),
    		xgen.WithProject(
    			"myproject",
    			basic.NewProject(),
    		),
    		xgen.WithProject(
    			"gormproj",
    			gorm.NewProject(
    				gorm.WithBasicProjectOption(basic.WithTestDir("tests")),
    			),
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
├── gormproj
│   ├── generated
│   │   ├── generated_gqlgen.go
│   │   ├── generated_gqlgen_models.go
│   │   └── generated_xgen_mappers.go
│   ├── tests
│   │   ├── default-tests.yaml
│   │   ├── user-lifecycle.yaml
│   │   └── user-pagination.yaml
│   ├── generated_xgen_directives.graphql
│   ├── generated_xgen_gorm.go
│   ├── generated_xgen_introspection.go
│   ├── generated_xgen_introspection.graphql
│   ├── generated_xgen_project_handlers.go
│   ├── graphql.config.yml
│   ├── resolver.go
│   ├── schema.main.graphql
│   └── schema.resolver.go
├── myproject
│   ├── generated
│   │   ├── generated_gqlgen.go
│   │   ├── generated_gqlgen_models.go
│   │   └── generated_xgen_mappers.go
│   ├── tests
│   │   └── default-tests.yaml
│   ├── generated_xgen_directives.graphql
│   ├── generated_xgen_introspection.go
│   ├── generated_xgen_introspection.graphql
│   ├── generated_xgen_project_handlers.go
│   ├── graphql.config.yml
│   ├── resolver.go
│   ├── schema.main.graphql
│   ├── schema.resolver.go
│   ├── schema.todo.graphql
│   └── schema.users.graphql
├── .gitignore
├── gen.go
├── generated_xgen_cli.go
├── gormproj.db
└── xgenc.go

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
	"github.com/goxgen/goxgen/utils/mapper"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Resolver struct {
	DB     *gorm.DB
	Mapper *mapper.Mapper
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}
	db, err := NewGormDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm db: %w", err)
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
    phoneNumbers: [Phone!]! @Field(Label: "Phone Numbers", Description: "Phone numbers of the user", DB: {})
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
@Action(Resource: "car", Action: CREATE_MUTATION, Route: "new")
@Action(Resource: "car", Action: UPDATE_MUTATION, Route: "update")
{
    id: ID @ActionField(Label: "ID", Description: "ID of the car", MapTo: ["Car.ID"])
    make: String @ActionField(Label: "Make", Description: "Text of the todo", MapTo: ["Car.Make"])
    done: Boolean @ActionField(Label: "Done", Description: "Done of the todo", MapTo: ["Car.Done"])
    user: UserInput @ActionField(Label: "User", Description: "User of the todo", MapTo: ["Car.User"])
}

input UserInput
@Action(Resource: "user", Action: CREATE_MUTATION, Route: "new")
@Action(Resource: "user", Action: UPDATE_MUTATION, Route: "update")
{
    id: ID @ActionField(Label: "ID", Description: "ID of the user", MapTo: ["User.ID"])
    name: String @ActionField(Label: "Name", Description: "Name", MapTo: ["User.Name"])
    cars: [CarInput!] @ActionField(Label: "Cars", Description: "Cars of the user", MapTo: ["User.Cars"])
    phones: [PhoneNumberInput!] @ActionField(Label: "Phone Numbers", Description: "Phone numbers of the user", MapTo: ["User.PhoneNumbers"])
}

input DeleteUsers
@ListAction(Resource: "user", Action: BATCH_DELETE_MUTATION, Route: "delete")
{
    ids: [ID!] @ActionField(Label: "IDs", Description: "IDs of users")
}

input ListUser
@ListAction(Resource: "user", Action: BROWSE_QUERY, Route: "list", Pagination: true)
{
    id: ID @ActionField(Label: "ID", Description: "ID", MapTo: ["User.ID"])
    name: String @ActionField(Label: "Name", Description: "Name", MapTo: ["User.Name"])
}

input CarBrowseInput
@ListAction(Resource: "car", Action: BROWSE_QUERY, Route: "list")
{
    id: ID @ActionField(Label: "ID", Description: "ID")
    userId: ID @ActionField(Label: "User ID", Description: "User ID")
    make: String @ActionField(Label: "Make", Description: "Make")
}

type Phone
@Resource(Name: "phone_number", Primary: true, Route: "phone-number", DB: {Table: "phone_number"})
{
    id: ID! @Field(Label: "ID", Description: "ID of the phone number", DB: {Column: "id", PrimaryKey: true})
    number: String! @Field(Label: "Number", Description: "Number of phone", DB: {Column: "number"})
    user: User! @Field(Label: "User", Description: "User of the todo", DB: {})
}

input PhoneNumberBrowseInput
@ListAction(Resource: "phone_number", Action: BROWSE_QUERY, Route: "list")
{
    id: ID @ActionField(Label: "ID", Description: "ID", MapTo: ["Phone.ID"])
    number: ID @ActionField(Label: "Number", Description: "Number of phone", MapTo: ["Phone.Number"])
}

input PhoneNumberInput
@Action(Resource: "phone_number", Action: CREATE_MUTATION, Route: "new")
@Action(Resource: "phone_number", Action: UPDATE_MUTATION, Route: "update")
{
    id: ID @ActionField(Label: "ID", Description: "ID of the phone number", MapTo: ["Phone.ID"])
    number: String @ActionField(Label: "Name", Description: "Number of phone", MapTo: ["Phone.Number"])
    user: UserInput @ActionField(Label: "User", Description: "User of the phone", MapTo: ["Phone.User"])
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
directive @ActionField(Label: String, Description: String,
  """Map field to resource field, {resource}.{field}, eg. user.id"""
  MapTo: [String!]
) on INPUT_FIELD_DEFINITION
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
  size: Int!
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
You can find the resolver functions for each field in the `schema.resolver.go` file.

### "Create User" mutation resolver
```go
func (r *mutationResolver) UserCreate(ctx context.Context, input *generated.UserInput) (*generated.User, error) {
	u, err := input.ToUserModel(ctx)
	if err != nil {
		return nil, err
	}
	res := r.DB.Preload(clause.Associations).Create(u)
	if res.Error != nil {
		return nil, res.Error
	}
	return u, nil
}
```

### "Browse User" query resolver
```go
func (r *queryResolver) UserBrowse(ctx context.Context, where *generated.ListUser, pagination *generated.XgenPaginationInput) ([]*generated.User, error) {
	// Get logger from context
	logger := server.GetLogger(ctx)
	logger.Info("UserBrowse", zap.Any("where", where))

	var users []*generated.User
	u, err := where.ToUserModel(ctx)
	if err != nil {
		return nil, err
	}
	res := r.DB.
		Preload(clause.Associations).
		Scopes(Paginate(pagination)).
		Where(&[]*generated.User{u}).
		Find(&users)

	if res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}
```
etc.

You can add your own implementation for each function in the updated `schema.resolver.go` file.
For more information,
You can read the [gqlgen documentation](https://gqlgen.com/getting-started/#implement-the-resolvers). 

In those functions, you can see that the `r.DB` instance is used, 
which is provided from the `resolver.go` file.

```go
package gormproj

import (
	"fmt"
	"github.com/goxgen/goxgen/utils/mapper"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Resolver struct {
	DB     *gorm.DB
	Mapper *mapper.Mapper
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}
	db, err := NewGormDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm db: %w", err)
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
    user1: user_create(input: {name: "My user 1", cars: {make:"BMW"}}){
        id
        name
        cars {
            make
            id
        }
    }
    user2: user_create(input: {name: "My user 2", cars: {make:"Mercedes"}}){
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
          "id": 1
        }
      ]
    },
    "user2": {
      "id": 2,
      "name": "My user 2",
      "cars": [
        {
          "make": "Mercedes",
          "id": 2
        }
      ]
    }
  }
}

```

One more example, let's list our new users by query

```graphql
query{
    user_browse(where: {}){
        id
        name
    }
}

```

The result of this query should be like this

```graphql
{
  "data": {
    "user_browse": [
      {
        "id": 1,
        "name": "My user 1"
      },
      {
        "id": 2,
        "name": "My user 2"
      }
    ]
  }
}

```

## 🤝 Contributing

> To configure git hooks, run `make init`

Contributions, issues, and feature requests are welcome!

### Makefile

To simplify the development process, we use Makefile.

- `make init` - Initialize git hooks
- `make pre-commit` - Run pre-commit checks
- `make integrations-generate` - Generate an integration test project
- `make integrations-run` - Run integration test project
- `make runtime-generate` - Generate a runtime project that using for goxgen code generation
- `make build-readme` - Build README.md file from README.gomd
- `make build` - Build all and prepare release

## 📝 License

Apache 2.0

## 📞 Contact

For more information, feel free to open an issue in the repository.

---

Enjoy the power of single-syntax API and domain definitions with `goxgen`! 🚀