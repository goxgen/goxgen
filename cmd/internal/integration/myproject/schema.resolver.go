package myproject

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/goxgen/goxgen/cmd/internal/integration/myproject/generated"
)

// NewTodo is the resolver for the new_todo field.
func (r *mutationResolver) NewTodo(ctx context.Context, input *generated.NewTodo) (*generated.Todo, error) {
	panic(fmt.Errorf("not implemented: NewTodo - new_todo"))
}

// TodoDelete is the resolver for the todo_delete field.
func (r *mutationResolver) TodoDelete(ctx context.Context, input *generated.DeleteTodo) (*generated.Todo, error) {
	return &generated.Todo{
		ID: input.ID,
	}, nil
}

// UserUpdate is the resolver for the user_update field.
func (r *mutationResolver) UserUpdate(ctx context.Context, input *generated.UpdateUser) (*generated.User, error) {
	return &generated.User{
		ID:   1,
		Name: input.Name,
	}, nil
}

// TodoCustom is the resolver for the todo_custom field.
func (r *mutationResolver) TodoCustom(ctx context.Context, input *generated.CustomTodo) (*generated.Todo, error) {
	return &generated.Todo{
		ID:   1,
		Text: input.Text,
		Done: false,
		User: &generated.User{
			ID: input.UserID,
		},
	}, nil
}

// TodoUpdate is the resolver for the todo_update field.
func (r *mutationResolver) TodoUpdate(ctx context.Context, input *generated.UpdateTodo) (*generated.Todo, error) {
	return &generated.Todo{
		ID:   1,
		Text: input.Text,
		Done: false,
		User: &generated.User{
			ID: input.UserID,
		},
	}, nil
}

// UserCreate is the resolver for the user_create field.
func (r *mutationResolver) UserCreate(ctx context.Context, input *generated.NewUser) (*generated.User, error) {
	return &generated.User{
		ID:   1,
		Name: input.Name,
	}, nil
}

// XgenIntrospection is the resolver for the _xgen_introspection field.
func (r *queryResolver) XgenIntrospection(ctx context.Context) (*generated.XgenIntrospection, error) {
	return r.Resolver.XgenIntrospection()
}

// TodoBrowse is the resolver for the todo_browse field.
func (r *queryResolver) TodoBrowse(ctx context.Context, input *generated.ListTodo, pagination *generated.XgenPaginationInput) ([]*generated.Todo, error) {
	return []*generated.Todo{
		{
			ID:   1,
			Text: input.Text,
			User: &generated.User{
				ID: input.UserID,
			},
		},
		{
			ID:   1,
			Text: input.Text,
			User: &generated.User{
				ID: input.UserID,
			},
		},
		{
			ID:   1,
			Text: input.Text,
			User: &generated.User{
				ID: input.UserID,
			},
		},
	}, nil
}

// UserBrowse is the resolver for the user_browse field.
func (r *queryResolver) UserBrowse(ctx context.Context, input *generated.UserList, pagination *generated.XgenPaginationInput) ([]*generated.User, error) {
	return []*generated.User{
		{
			ID:   1,
			Name: "John",
		},
		{
			ID:   2,
			Name: "Doe",
		},
		{
			ID:   3,
			Name: "Smith",
		},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
