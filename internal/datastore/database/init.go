package database

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

// @Injectable
func InitDatabase() (*bun.DB, error) {
	sqldb, err := sql.Open(
		sqliteshim.ShimName,
		"file:test.db?cache=shared&mode=rwc",
	)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	return db, nil
}
