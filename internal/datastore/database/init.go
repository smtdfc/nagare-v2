package database

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/smtdfc/nagare/core/datastore/database/models"
	"github.com/smtdfc/nagare/core/datastore/filesystem"
	error_utils "github.com/smtdfc/nagare/core/error-utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func RunMigration(db *bun.DB) error {
	ctx := context.Background()
	models := []interface{}{
		(*models.ModelProvider)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// @Injectable
func InitDatabase(paths *filesystem.DataPaths) (*bun.DB, error) {
	if paths == nil || paths.MainDatabaseFile == "" {
		return nil, error_utils.ErrDatabasePathEmpty
	}

	dsn := url.URL{
		Scheme:   "file",
		Path:     paths.MainDatabaseFile,
		RawQuery: "cache=shared&mode=rwc",
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn.String())
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	err = RunMigration(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}
