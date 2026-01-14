package storage

import (
	"log"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	db := NewDatabase()

	table, err := db.CreateTable("")

	if err != nil {
		t.Fatalf("an error occurred when creating a table:%s", err)
	}

	if table.Name != "users" { // Assertion
		t.Fatalf("expected table name 'users', got %v", table.Name)
	}

	if len(db.Tables) != 0 {
		log.Fatalf("expected empty tables map, got %d tables", len(db.Tables))
	}
}
