package myproject

import (
	"github.com/urfave/cli/v2"
)

type Resolver struct{}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	return &Resolver{}, nil
}
