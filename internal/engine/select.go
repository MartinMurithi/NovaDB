package engine

import (
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

// GetByPK retrieves a single row by its primary key from the specified table.
func (e *Engine) GetByPK(tableName string, pk any) (*storage.Row, error) {
	// checks if table exists
	table, ok := e.db.Tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	return table.GetRowByPK(pk)
}

// SelectAll returns all rows in the specified table.
func (e *Engine) SelectAll(tableName string) ([]*storage.Row, error) {
	table, ok := e.db.Tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	return table.Rows, nil
}

// SelectByColumnValue returns all rows in a table where the given column matches a value.
func (e *Engine) SelectByColumnValue(tableName, columnName string, value any) ([]*storage.Row, error) {
	table, ok := e.db.Tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	var result []*storage.Row
	for _, row := range table.Rows {
		if rowVal, exists := row.Data[columnName]; exists && rowVal == value {
			result = append(result, row)
		}
	}

	return result, nil
}