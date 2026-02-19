package testhelper

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// OpenTestDB returns a database connection for integration tests.
// Skips the test if TEST_DATABASE_URL is not set, so unit tests are unaffected.
func OpenTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping integration test")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("ping test db: %v", err)
	}

	t.Cleanup(func() {
		truncateAll(t, db)
		db.Close()
	})

	return db
}

func truncateAll(t *testing.T, db *sql.DB) {
	t.Helper()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		t.Logf("show tables: %v", err)
		return
	}
	defer rows.Close()

	// Disable FK checks for safe truncation
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		t.Logf("disable fk checks: %v", err)
		return
	}

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			t.Logf("scan table name: %v", err)
			continue
		}
		// Keep migration tracking table intact
		if table == "schema_migrations" {
			continue
		}
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)); err != nil {
			t.Logf("truncate %s: %v", table, err)
		}
	}

	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		t.Logf("enable fk checks: %v", err)
	}
}
