package gormproj

import (
	"fmt"
	"github.com/goxgen/goxgen/utils/mapper"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Resolver struct {
	DB     *gorm.DB
	Mapper *mapper.Mapper
}

func NewResolver(ctx *cli.Context) (*Resolver, error) {
	r := &Resolver{}
	db, err := NewGormDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm db: %w", err)
	}
	r.DB = db
	return r, nil
}
