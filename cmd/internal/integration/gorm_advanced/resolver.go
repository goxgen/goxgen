package gorm_advanced

import (
	"github.com/goxgen/goxgen/cmd/internal/integration/gorm_advanced/generated"
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