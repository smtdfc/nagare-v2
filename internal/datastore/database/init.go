package database

import (
	"database/sql"
	"errors"
	"net/url"

	"github.com/smtdfc/nagare/core/datastore/filesystem"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

var ErrDatabasePathEmpty = errors.New("database: main database file path is empty")

// @Injectable
func InitDatabase(paths *filesystem.DataPaths) (*bun.DB, error) {
	if paths == nil || paths.MainDatabaseFile == "" {
		return nil, ErrDatabasePathEmpty
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
	return db, nil
}
