package simple_initial

import (
	"github.com/goxgen/goxgen/plugins/cli/settings"
)


type Resolver struct{}

func NewResolver(sts *settings.EnvironmentSettings) (*Resolver, error) {
	return &Resolver{}, nil
}