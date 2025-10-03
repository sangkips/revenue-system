package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

func ConnectAndMigrate(dbURL string) *sql.DB {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")
	}

	log.Info().Msg("Successfully connected to database")
	return db
}

func SetUpTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dbURL := "postgres://test_user:test_password@localhost:5435/test_county_db?sslmode=disable"

	if url := os.Getenv("TEST_DB_URL"); url != "" {
		dbURL = url
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")

	}

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("Failed to ping database")

	}
	log.Info().Msg("Successfully connected to test database")

	return db
}

func TeardownTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec("TRUNCATE TABLE users, counties CASCADE"); err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
	db.Close()
}
