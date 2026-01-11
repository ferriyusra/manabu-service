// +build ignore

package main

import (
	"fmt"
	"log"
	"manabu-service/config"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/migrate.go <migration_file>")
		fmt.Println("Example: go run tools/migrate.go 002_rename_users_uuid_constraint.sql")
		os.Exit(1)
	}

	migrationFile := os.Args[1]

	// Initialize config
	config.Init()

	// Connect to database
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Read migration file
	migrationPath := filepath.Join("migrations", migrationFile)
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	sqlContent := string(sqlBytes)

	// Execute migration
	fmt.Printf("Running migration: %s\n", migrationFile)
	fmt.Println("---")
	fmt.Println(sqlContent)
	fmt.Println("---")
	fmt.Print("\nContinue? (yes/no): ")

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		fmt.Println("Migration cancelled.")
		return
	}

	// Execute SQL
	result := db.Exec(sqlContent)
	if result.Error != nil {
		log.Fatalf("Migration failed: %v", result.Error)
	}

	fmt.Println("\n✓ Migration completed successfully!")
}
