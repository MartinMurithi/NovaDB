package engine

import "fmt"

// Delete removes a row by primary key from the specified table.
func (e *Engine) Delete(tableName string, pk any) error {
	table, ok := e.db.Tables[tableName]
	if !ok {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	return table.Delete(pk)
}