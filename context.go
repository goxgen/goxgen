package goxgen

import (
	"context"
	"fmt"
)

// Context is a context for Xgen
// it is used to pass data between generators
type Context struct {
	ParentPackageName string
	Projects          []Project
}

// GetContext gets generator context from context
// it is used to pass data between generators
func GetContext(ctx context.Context) (*Context, error) {
	gCtx, ok := ctx.Value(GeneratorContextKey).(Context)

	if !ok {
		return nil, fmt.Errorf("failed to get generator context")
	}
	return &gCtx, nil
}
