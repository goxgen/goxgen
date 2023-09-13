package gorm_initial

import (
    "github.com/goxgen/goxgen/runtime/gorm_initial/generated"
    "github.com/urfave/cli/v2"
    "gorm.io/driver/sqlite"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlserver"
    "gorm.io/driver/clickhouse"
    "gorm.io/gorm"
	"fmt"
)

func NewGormDB(ctx *cli.Context) (*gorm.DB, error){
    dbDriver := ctx.String("DatabaseDriver")
    dbDsn := ctx.String("DatabaseSourceName")
    var dialector gorm.Dialector
    switch dbDriver {
    case "sqlite":
        dialector = sqlite.Open(dbDsn)
    case "mysql":
        dialector = mysql.Open(dbDsn)
    case "postgres":
        dialector = postgres.Open(dbDsn)
    case "sqlserver":
        dialector = sqlserver.Open(dbDsn)
    case "clickhouse":
        dialector = clickhouse.Open(dbDsn)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
    }

    // Open the database connection
    db, err := gorm.Open(
        dialector,
        &gorm.Config{
			FullSaveAssociations: true,
        },
    )
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Migrate the schema
    err = db.AutoMigrate(
    )
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    return db, nil
}