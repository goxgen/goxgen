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
    //go:build ignore
    // +build ignore
    
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
    		xgen.WithPackageName("github.com/goxgen/goxgen/cmd/internal/integration"),
    		xgen.WithProject(
    			"myproject",
    			basic.NewProject(),
    		),
    		xgen.WithProject(
    			"gorm_advanced",
    			gorm.NewProject(
    				gorm.WithBasicProjectOption(basic.WithTestDir("tests")),
    			),
    		),
    		xgen.WithProject(
    			"gorm_example",
    			gorm.NewProject(
    				gorm.WithBasicProjectOption(basic.WithTestDir("tests")),
    			),
    		),
    		xgen.WithPlugin(cli.NewPlugin()),
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
├── gorm_advanced
│   ├── generated
│   │   ├── server
│   │   ├── generated_gqlgen.go
│   │   ├── generated_gqlgen_models.go
│   │   ├── generated_xgen_directives.graphql
│   │   ├── generated_xgen_gorm.go
│   │   ├── generated_xgen_introspection.go
│   │   ├── generated_xgen_introspection.graphql
│   │   ├── generated_xgen_mappers.go
│   │   └── generated_xgen_sortable.go
│   ├── tests
│   │   ├── default-tests.yaml
│   │   ├── user-lifecycle.yaml
│   │   └── user-pagination.yaml
│   ├── graphql.config.yml
│   ├── resolver.go
│   ├── schema.main.graphql
│   └── schema.resolver.go
├── gorm_example
│   ├── generated
│   │   ├── server
│   │   ├── generated_gqlgen.go
│   │   ├── generated_gqlgen_models.go
│   │   ├── generated_xgen_directives.graphql
│   │   ├── generated_xgen_gorm.go
│   │   ├── generated_xgen_introspection.go
│   │   ├── generated_xgen_introspection.graphql
│   │   ├── generated_xgen_mappers.go
│   │   └── generated_xgen_sortable.go
│   ├── tests
│   │   ├── default-tests.yaml
│   │   └── user-lifecycle.yaml
│   ├── graphql.config.yml
│   ├── resolver.go
│   ├── schema.phone.graphql
│   ├── schema.resolver.go
│   └── schema.user.graphql
├── myproject
│   ├── generated
│   │   ├── server
│   │   ├── generated_gqlgen.go
│   │   ├── generated_gqlgen_models.go
│   │   ├── generated_xgen_directives.graphql
│   │   ├── generated_xgen_introspection.go
│   │   ├── generated_xgen_introspection.graphql
│   │   ├── generated_xgen_mappers.go
│   │   └── generated_xgen_sortable.go
│   ├── tests
│   │   └── default-tests.yaml
│   ├── graphql.config.yml
│   ├── resolver.go
│   ├── schema.main.graphql
│   ├── schema.resolver.go
│   ├── schema.todo.graphql
│   └── schema.users.graphql
├── .env
├── .env.default
├── .gitignore
├── gen.go
├── generated_xgen_cli.go
├── gorm_advanced.db
├── gorm_example.db
├── gormproj.db
└── xgenc.go

```

> Note: `generated` directories can be ignored in git. But you can add it to git if you want.

### 📑 Providing schema

You should provide a schema for each project and run `go generate` again.

All schema files in xgen has this format `schema.{some_name}.graphql`, for example `schema.user.graphql`

#### Gorm example

Let's focus on `gorm_example`, which uses the GORM ORM.
The connection to the GORM database can be configured from the gqlgen standard `resolver.go` file in the `gorm_example` directory.

> `resolver.go` is designed to support your custom dependency injection (DI) and any services you've provided.

```go
package gorm_example

import (
	"github.com/goxgen/goxgen/cmd/internal/integration/gorm_example/generated"
	"github.com/goxgen/goxgen/plugins/cli/settings"
	"gorm.io/gorm"
	"embed"
	"fmt"
)

//go:embed tests/*
var TestsFS embed.FS

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(sts *settings.EnvironmentSettings) (*Resolver, error) {
	r := &Resolver{}
	db, err := generated.NewGormDB(sts)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm db: %w", err)
	}
	r.DB = db

	return r, nil
}
```

### Example of schema file `schema.user.graphql`

```graphql
# Define the User resource(entity) and its fields
# Enable DB mapping for the resource
type User
@Resource(Name: "user", Primary: true, Route: "user", DB: {Table: "user"})
{
    id: ID! @Field(Label: "ID", Description: "ID of the user", DB: {Column: "id", PrimaryKey: true})
    name: String! @Field(Label: "Text", Description: "Text of the user", DB: {Column: "name", Unique: true})
    phoneNumbers: [Phone!]! @Field(Label: "Phone Numbers", Description: "Phone numbers of the user", DB: {})
}

# User input type for create and update actions
# Define the actions for the resource
input UserInput
@Action(Resource: "user", Action: CREATE_MUTATION, Route: "new")
@Action(Resource: "user", Action: UPDATE_MUTATION, Route: "update")
{
    id: ID @ActionField(Label: "ID", Description: "ID of the user", MapTo: ["User.ID"])
    name: String @ActionField(Label: "Name", Description: "Name", MapTo: ["User.Name"])
    phones: [PhoneNumberInput!] @ActionField(Label: "Phone Numbers", Description: "Phone numbers of the user", MapTo: ["User.PhoneNumbers"])
}

# User input type for browse action
input BrowseUserInput
@ListAction(Resource: "user", Action: BROWSE_QUERY, Route: "list", Pagination: true, Sort: {Default: [{by: "name", direction: ASC}]})
{
    id: ID @ActionField(Label: "ID", Description: "ID", MapTo: ["User.ID"])
    name: String @ActionField(Label: "Name", Description: "Name", MapTo: ["User.Name"])
}
```

The directives used in the example above are standard `xgen` directives, intended to provide metadata.

* `Resource` - Entity or object or thing
* `Field` - Field of resource
* `Action` - Action that can be done for single resource
* `ListAction` - Action that can be done for bulk resources
* `ActionField` - Field of action or list action

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
func (r *queryResolver) UserBrowse(ctx context.Context, where *generated.BrowseUserInput, pagination *generated.XgenPaginationInput, sort *generated.UserSortInput) ([]*generated.User, error) {
	var users []*generated.User
	u, err := where.ToUserModel(ctx)
	if err != nil {
		return nil, err
	}
	res := r.DB.
		Preload(clause.Associations).
		Scopes(
			generated.Paginate(pagination), // passing `pagination` to the xgen `generated.Paginate` scope
			generated.Sort(sort),           // passing `sort` to the xgen `generated.Sort` scope
		).
		Where(&[]*generated.User{u}).
		Find(&users)

	return users, res.Error
}
```
etc.

You can add your own implementation for each function in the updated `schema.resolver.go` file.
For more information,
You can read the [gqlgen documentation](https://gqlgen.com/getting-started/#implement-the-resolvers). 

In those functions, you can see that the `r.DB` instance is used, 
which is provided from the `resolver.go` file.

```go
package gorm_example

import (
	"github.com/goxgen/goxgen/cmd/internal/integration/gorm_example/generated"
	"github.com/goxgen/goxgen/plugins/cli/settings"
	"gorm.io/gorm"
	"embed"
	"fmt"
)

//go:embed tests/*
var TestsFS embed.FS

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(sts *settings.EnvironmentSettings) (*Resolver, error) {
	r := &Resolver{}
	db, err := generated.NewGormDB(sts)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm db: %w", err)
	}
	r.DB = db

	return r, nil
}
```

Great, you're all set to launch your GraphQL application.

### 🖥️ CLI plugin usage

To start the server using the xgen CLI plugin, you can run the following command:

```shell
go run generated_xgen_cli.go run --gql-playground-enabled
```

This will initialize and start your all projects GraphQL servers together, making it ready to handle incoming requests.

The output from the xgen CLI will provide information about the server endpoints. Additionally, logs will be written to this output during the server's runtime, giving you insights into its operation.

```shell
2023-10-09T00:46:43.600+0400    INFO    server/server.go:77     Serving graphql playground      {"project": "gorm_example", "url": "http://localhost:8080/playground"}
2023-10-09T00:46:43.600+0400    INFO    server/server.go:88     Serving graphql                 {"project": "gorm_example", "url": "http://localhost:8080/query"}
```

If You have a more then one project, and you want to run only one or some projects, you can use `--project` flag

```shell
go run generated_xgen_cli.go run --gql-playground-enabled --project gorm_example
```

Or for multiple projects
```shell
go run generated_xgen_cli.go run --gql-playground-enabled --project gorm_example --project otherproj
```

#### 📖 GraphQL Playground
To enable the GraphQL playground, you can use the `--gql-playground-enabled` flag.

#### 🔡 Environment variables
By default, the xgen generating two dotenv files in your root directory - `.env` and `.env.default`.

* `.env.default` file is auto-generated and contains necessary environment variables for your project. Do not edit this file because it will be overwritten on each generation.
    ```properties
    # Auto generated by goxgen, do not edit manually
    # This is default environment variables for github.com/goxgen/goxgen/cmd/internal/integration project
    
    # gorm_advanced project default environment variables
    GORM_ADVANCED_PORT=8080
    GORM_ADVANCED_DB_DRIVER=sqlite
    GORM_ADVANCED_DB_DSN=file:gorm_advanced.db?mode=rwc&cache=shared&_fk=1
    
    # gorm_example project default environment variables
    GORM_EXAMPLE_PORT=8081
    GORM_EXAMPLE_DB_DRIVER=sqlite
    GORM_EXAMPLE_DB_DSN=file:gorm_example.db?mode=rwc&cache=shared&_fk=1
    
    # myproject project default environment variables
    MYPROJECT_PORT=8082
    MYPROJECT_DB_DRIVER=sqlite
    MYPROJECT_DB_DSN=file:myproject.db?mode=rwc&cache=shared&_fk=1
    ```
* `.env` file is a file that you can edit and add your own environment variables. This file is not overwritten on each generation.

You can also use .env.local file for local environment variables.

##### Structure of environment variables
Xgen CLI has a special structure for environment variables.
You can define default environment variables for all projects
and override them for each project with project name prefix.

```properties
{ENVIRONMENT_VARIABLE_NAME}={VALUE}
{PROJECT_NAME}_{ENVIRONMENT_VARIABLE_NAME}={VALUE}
```

e.g.
```properties
# Default environment variable for all projects
DB_DSN=sqllite://file.db
# Environment variable for gorm_example project
GORM_EXAMPLE_DB_DSN=postgres://user:pass@localhost:5432/gorm_example?sslmode=disable
```

##### Available environment variables
To see all available environment variables, you can run the following command:

```shell
go run generated_xgen_cli.go run --help
```

> For more information about the xgen CLI, you can run main help command:
> 
> `go run generated_xgen_cli.go help`
> 
> This will display a list of available commands, options, and descriptions to help you navigate the xgen CLI more effectively.

## Playground and testing
You can copy the URL `http://localhost:80/playground` from the logs 
and open it in your browser to access the GraphQL playground. 
This interface will allow you to test queries, mutations, and subscriptions in real-time.

Then we see graphql playground, let's run some mutation query to add two new users

```graphql
mutation{
  user1: user_create(input: {name: "My user 1"}){
      id
      name
  }
  user2: user_create(input: {name: "My user 2"}){
      id
      name
  }
}

```

After execution of this mutation, graphql should be return result like this

```json
{
    "user1": {
      "id": 1,
      "name": "My user 1"
    },
    "user2": {
      "id": 2,
      "name": "My user 2"
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

```

## Testing

Xgen has a support for custom api tests. You can write your own tests in yaml format and run it CLI command.

In generated project directory you can find `tests` directory. Xgen also generates a default test file `tests/default-tests.yaml`.

```yaml
name: "Default tests"
tests:
    - name: "Healthcheck"
      query: |
        query{
          __schema{
            __typename
          }
        }
      expectedResult: |
        {
          "__schema": {
            "__typename": "__Schema"
          }
        }
```

You can create your own test file and run it with CLI command.

```shell
go run generated_xgen_cli.go run --test
```

This command will run all tests in all projects. If you want to run tests only for one project, you can use `--project` flag.

## Available Project Types

### Basic Project
Basic project is a project without any ORM. It's a simple project with a simple structure.
You can use it for your own custom implementation.

### Gorm Project
Gorm project is a project with GORM ORM.

#### Pagination and Sorting
Resolver method `UserBrowse` has a `Pagination` and `Sort` arguments. This arguments is a set of standard pagination and sort parameters.
Xgen provides a special GORM scopes for pagination and sort functionalities. You can use it in your custom implementation.

```go
func (r *queryResolver) UserBrowse(ctx context.Context, where *generated.BrowseUserInput, pagination *generated.XgenPaginationInput, sort *generated.UserSortInput) ([]*generated.User, error) {
	var users []*generated.User
	u, err := where.ToUserModel(ctx)
	if err != nil {
		return nil, err
	}
	res := r.DB.
		Preload(clause.Associations).
		Scopes(
			generated.Paginate(pagination), // passing `pagination` to the xgen `generated.Paginate` scope
			generated.Sort(sort),           // passing `sort` to the xgen `generated.Sort` scope
		).
		Where(&[]*generated.User{u}).
		Find(&users)

	return users, res.Error
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