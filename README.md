# goxgen

Library for generating graphql application based on domain definitions.
> One syntax for defining domain and api interface.

Using graphql as a schema language, and generating graphql application (server) for this application, 
providing ORM support and CLI to run server application. 
In the future, planning to provide a UI for admin backoffice and authorization and authentication.

> goxgen based on gqlgen, and using it as a base for generating graphql application

# How it works

## Step-by-step guide

### Creating the necessary files

You should create two files in your project

1. Standard `gen.go` file with `go:generate` directive
    ```go
    package main
    
    //go:generate go run -mod=mod github.com/goxgen/goxgen
    
    ```
2. Xgen config file `xgenc.go`
   https://github.com/goxgen/goxgen-demo/blob/main/xgenc.go#L12-L34

Then run `go generate` command, and goxgen will generate project structure

```shell
go generate
```

### Structure of a generated project

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

You should provide a schema for each project and run `go generate` again.

All schema files in xgen has this format `schema.{some_name}.graphql`, for example `schema.main.graphql`


Let's focus on `gormproj`, which uses the GORM ORM.
The connection to the GORM database can be configured from the gqlgen standard `resolver.go` file in the `gormproj` directory.

> `resolver.go` is designed to support your custom dependency injection (DI) and any services you've provided.

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/resolver.go#L11-L35

### Example of schema file `schema.main.graphql`

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/schema.main.graphql

The directives used in the example above are standard `xgen` directives, intended to provide metadata.

* `Resource` - Entity or object or thing
* `Field` - Field of resource
* `Action` - Action that can be done for single resource
* `ListAction` - Action that can be done for bulk resources
* `ActionField` - Field of action or list action

The definitions of these directives are located in the `generated_xgen_directives.graphql` file.

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/generated_xgen_directives.graphql#L1-L12

After writing a custom schema You should run again `gogen` command.

```shell
go generate
```

After regenerating the code, the `schema.resolver.go` file will be updated based on your schema. 
You can expect to see changes similar to the following:

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/schema.resolver.go#L16-L90

You can add your own implementation for each function in the updated `schema.resolver.go` file.
For more information,
You can read the [gqlgen documentation](https://gqlgen.com/getting-started/#implement-the-resolvers). 


You can see that some functions in example are implemented, some are not. 
So take a look at the implemented functions, for example `NewUser`, `ListUser` or `DeleteUsers`.
In those functions, you can see that the `r.DB` instance is used, 
which is provided from the `resolver.go` file.

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/resolver.go#L11-L13

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

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/examples/new_user_mutation.graphql

After execution of this mutation, graphql should be return result like this
https://github.com/goxgen/goxgen-demo/blob/main/gormproj/examples/new_user_mutation_result.json

One more example, let's list our new users by query
https://github.com/goxgen/goxgen-demo/blob/main/gormproj/examples/list_user_query.graphql

The result of this query should be like this

https://github.com/goxgen/goxgen-demo/blob/main/gormproj/examples/list_user_query_result.json