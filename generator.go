package goxgen

import "context"

type Generator interface {
	Generate(ctx context.Context) error
}
