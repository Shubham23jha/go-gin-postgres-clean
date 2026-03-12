package integration

import (
	"log"
	"os"
	"testing"

	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDBURL = "postgres://user:password@localhost:5433/testdb?sslmode=disable"

func cleanupDB(db *gorm.DB) {
	tables := []string{"user_sessions", "users"}
	for _, table := range tables {
		db.Exec(`TRUNCATE TABLE "` + table + `" CASCADE`)
	}
	log.Println("🧹 Test database cleaned")
}

// TestMain handles the setup and teardown for all integration tests.
// DB Lifecycle:
// 1. Reset: The schema is completely dropped and recreated on every run (avoids dirty states).
// 2. Migrate: migrations are applied fresh.
// 3. Clean: Truncates data tables before the suite starts.
// 4. Persist: Data remains in the DB after tests finish to allow manual inspection/debugging.
func TestMain(m *testing.M) {
	// Set environment variables for the test environment
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("REFRESH_SECRET", "test-refresh-secret")

	// Run migrations
	migrationPath := "../../"
	if _, err := os.Stat("../../migrations"); os.IsNotExist(err) {
		migrationPath = "./"
	}

	// Initialize global DB connection for total reset
	dbReset, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err == nil {
		// Drop everything to fix dirty states
		dbReset.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
		log.Println("♻️ Database schema reset")
	}

	mSub, err := migrate.New(
		"file://"+migrationPath+"migrations",
		testDBURL,
	)
	if err != nil {
		log.Fatalf("failed to create migration instance: %s", err)
	}

	if err := mSub.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("failed to run migrations: %s", err)
	}
	log.Println("✅ Migrations applied to test database")

	// Initialize global DB connection for tests
	db, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test db: %s", err)
	}
	database.DB = db

	// Clean database before running tests
	cleanupDB(database.DB)

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}
