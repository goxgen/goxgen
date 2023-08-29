package simple_initial

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/goxgen/goxgen/runtime/simple_initial/generated_gqlgen"
)

// XgenIntrospection is the resolver for the _xgen_introspection field.
func (r *queryResolver) XgenIntrospection(ctx context.Context) (*generated_gqlgen.XgenIntrospection, error) {
	return r.Resolver.XgenIntrospection()
}

// Annotation is the resolver for the annotation field.
func (r *queryResolver) Annotation(ctx context.Context) (*generated_gqlgen.XgenAnnotationMap, error) {
	panic(fmt.Errorf("not implemented: Annotation - annotation"))
}

// Object is the resolver for the object field.
func (r *queryResolver) Object(ctx context.Context) (*generated_gqlgen.XgenObjectMap, error) {
	panic(fmt.Errorf("not implemented: Object - object"))
}

// Query returns generated_gqlgen.QueryResolver implementation.
func (r *Resolver) Query() generated_gqlgen.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }