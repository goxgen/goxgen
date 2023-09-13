package generated

import (
	"context"
	"github.com/goxgen/goxgen/plugins/cli/server"
)

// ToTodoModel Map DeleteTodo to Todo model
func (ra *DeleteTodo) ToTodoModel(ctx context.Context) (*Todo, error) {
	mapper := server.GetMapper(ctx)
	target := &Todo{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToTodoModel Map ListTodo to Todo model
func (ra *ListTodo) ToTodoModel(ctx context.Context) (*Todo, error) {
	mapper := server.GetMapper(ctx)
	target := &Todo{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToTodoModel Map NewTodo to Todo model
func (ra *NewTodo) ToTodoModel(ctx context.Context) (*Todo, error) {
	mapper := server.GetMapper(ctx)
	target := &Todo{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToTodoModel Map CustomTodo to Todo model
func (ra *CustomTodo) ToTodoModel(ctx context.Context) (*Todo, error) {
	mapper := server.GetMapper(ctx)
	target := &Todo{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToTodoModel Map UpdateTodo to Todo model
func (ra *UpdateTodo) ToTodoModel(ctx context.Context) (*Todo, error) {
	mapper := server.GetMapper(ctx)
	target := &Todo{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map UserList to User model
func (ra *UserList) ToUserModel(ctx context.Context) (*User, error) {
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map UpdateUser to User model
func (ra *UpdateUser) ToUserModel(ctx context.Context) (*User, error) {
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}

// ToUserModel Map NewUser to User model
func (ra *NewUser) ToUserModel(ctx context.Context) (*User, error) {
	mapper := server.GetMapper(ctx)
	target := &User{}
	err := mapper.Map(ra, target)
	return target, err
}
