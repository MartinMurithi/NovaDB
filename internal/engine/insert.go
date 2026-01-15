package engine

import (
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

// Insert inserts a new row into a table
func (e *Engine) Insert(tableName string, data map[string]any) error {
	if tableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	table, exists := e.db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	row := &storage.Row{
		Data: data,
	}

	return table.Insert(row)
}