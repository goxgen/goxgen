package gormproj

import (
	"fmt"
	"github.com/goxgen/goxgen/cmd/internal/integration/gormproj/generated"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}

	// Open the database connection
	db, err := gorm.Open(sqlite.Open("./cmd/internal/integration/gormproj.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&generated.Car{},
		&generated.User{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	r.DB = db
	return r, nil
}
