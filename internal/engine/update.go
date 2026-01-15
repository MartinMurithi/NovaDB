package engine

import "fmt"

// Update updates the values of a row identified by its primary key.
func (e *Engine) Update(tableName string, pk any, values map[string]any) error {
	table, ok := e.db.Tables[tableName]
	if !ok {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	return table.Update(pk, values)
}