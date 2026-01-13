package storage

import "fmt"

type Database struct {
	Tables map[string]*Table
}

// NewDatabase initializes and returns a new Database instance.
// The database starts with an empty table catalog.
func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*Table),
	}
}

// CreateTable creates a new empty table and registers it in the database.
//
// The table is created with no columns and no rows, but all internal
// structures are initialized and safe for use.
//
// Returns an error if the table name is empty or already exists.
func (db *Database) CreateTable(name string) (*Table, error) {
	// Ensure the table catalog is initialized
	if db.Tables == nil {
		db.Tables = make(map[string]*Table)
	}

	// Validate table name
	if name == "" {
		return nil, fmt.Errorf("table name cannot be empty")
	}

	// Enforce unique table names
	if _, exists := db.Tables[name]; exists {
		return nil, fmt.Errorf("table %s already exists", name)
	}

	// Initialize an empty table
	t := &Table{
		Name:         name,
		Columns:      make([]*Column, 0),
		Rows:         make([]*Row, 0),
		PrimaryIndex: make(map[any]int),
	}

	// Register the table in the database
	db.Tables[name] = t

	return t, nil
}
