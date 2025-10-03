package db

import (
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestDB_Connection(t *testing.T) {
	db := SetUpTestDB(t)
	defer TeardownTestDB(t, db)
}
