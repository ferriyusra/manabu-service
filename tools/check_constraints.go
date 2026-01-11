// +build ignore

package main

import (
	"fmt"
	"log"
	"manabu-service/config"
)

func main() {
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

	// Check constraints on users table
	fmt.Println("Checking constraints on 'users' table...")

	type Constraint struct {
		Name string
		Type string
	}

	var constraints []Constraint
	result := db.Raw(`
		SELECT conname as name, contype as type
		FROM pg_constraint
		WHERE conrelid = 'users'::regclass
		ORDER BY conname
	`).Scan(&constraints)

	if result.Error != nil {
		log.Fatalf("Failed to query constraints: %v", result.Error)
	}

	fmt.Println("\nConstraints found:")
	fmt.Println("------------------")
	for _, c := range constraints {
		typeDesc := ""
		switch c.Type {
		case "p":
			typeDesc = "PRIMARY KEY"
		case "u":
			typeDesc = "UNIQUE"
		case "f":
			typeDesc = "FOREIGN KEY"
		case "c":
			typeDesc = "CHECK"
		}
		fmt.Printf("- %s (%s)\n", c.Name, typeDesc)
	}

	// Check if uni_users_uuid exists
	hasUniUsersUUID := false
	for _, c := range constraints {
		if c.Name == "uni_users_uuid" {
			hasUniUsersUUID = true
			break
		}
	}

	fmt.Println("\nStatus:")
	if hasUniUsersUUID {
		fmt.Println("✓ Constraint 'uni_users_uuid' EXISTS (migration already applied)")
	} else {
		fmt.Println("✗ Constraint 'uni_users_uuid' NOT FOUND (migration needed)")
	}
}
