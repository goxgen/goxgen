package myproject

import (
	"embed"
	"github.com/goxgen/goxgen/plugins/cli/settings"
)

//go:embed tests/*
var TestsFS embed.FS

type Resolver struct{}

func NewResolver(sts *settings.EnvironmentSettings) (*Resolver, error) {
	return &Resolver{}, nil
}
