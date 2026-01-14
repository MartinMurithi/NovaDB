package storage

import "fmt"

// Table represents a database table with a name, schema (columns),
// row data, and a primary key index.
type Table struct {
	Name         string
	Columns      []*Column
	Rows         []*Row
	PrimaryIndex map[any]int // Maps primary key values to row indices
}

// AddColumn adds a new column to the table schema.
//
// Returns an error if the column is nil, has an empty name,
// or already exists in the table.
func (t *Table) AddColumn(c *Column) error {
	if c == nil {
		return fmt.Errorf("column cannot be empty")
	}

	if c.Name == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	for _, col := range t.Columns {
		if c.Name == col.Name {
			return fmt.Errorf("column %s already exists", c.Name)
		}
	}

	t.Columns = append(t.Columns, c)
	return nil
}

// DropColumn removes a column from the table schema by name.
//
// Returns an error if the column name is empty or does not exist.
// All rows and indexes related to the column should be updated separately.
func (t *Table) DropColumn(name string) error {
	if name == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	for index, col := range t.Columns {
		if col.Name == name {
			t.Columns = append(t.Columns[:index], t.Columns[index+1:]...)
			return nil
		}
	}

	return fmt.Errorf("column %s does not exist", name)
}

func (t *Table) Insert(row *Row) error {
	if row == nil {
		return fmt.Errorf("row cannot be nil")
	}

	// Check all columns in table exist in row
	for _, col := range t.Columns {
		if _, ok := row.Data[col.Name]; !ok {
			return fmt.Errorf("missing value for column %s", col.Name)
		}
	}

	// Check primary key uniqueness
	for _, col := range t.Columns {
		if col.IsPrimaryKey {
			val := row.Data[col.Name]
			if _, exists := t.PrimaryIndex[val]; exists {
				return fmt.Errorf("duplicate primary key value %v", val)
			}
			// Store index for quick lookup
			t.PrimaryIndex[val] = len(t.Rows)
			break
		}
	}

	// Add row
	t.Rows = append(t.Rows, row)
	return nil
}

// GetRows returns all rows in the table
func (t *Table) GetRows() []*Row {
	return t.Rows
}

// GetRowByPK retrieves a row by primary key value
func (t *Table) GetRowByPK(pk any) (*Row, error) {

	if t.PrimaryIndex == nil {
		return nil, fmt.Errorf("table has no primary key index")
	}

	index, exists := t.PrimaryIndex[pk]
	if !exists {
		return nil, fmt.Errorf("row with primary key %v not found", pk)
	}

	return t.Rows[index], nil
}

// FilterRows returns all rows matching a column-value pair
func (t *Table) FilterRows(column string, value any) ([]*Row, error) {
	var result []*Row

	// Check column exists
	found := false
	for _, col := range t.Columns {
		if col.Name == column {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("column %s does not exist", column)
	}

	for _, row := range t.Rows {
		if v, ok := row.Data[column]; ok && v == value {
			result = append(result, row)
		}
	}

	return result, nil
}

// Update updates a row identified by its primary key
func (t *Table) Update(pk any, updates map[string]any) error {
	if t.PrimaryIndex == nil {
		return fmt.Errorf("table has no primary key index")
	}

	index, exists := t.PrimaryIndex[pk]
	if !exists {
		return fmt.Errorf("row with primary key %v not found", pk)
	}

	row := t.Rows[index]

	// Update each column
	for colName, newValue := range updates {
		found := false
		for _, col := range t.Columns {
			if col.Name == colName {
				row.Data[colName] = newValue
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("column %s does not exist", colName)
		}
	}

	return nil
}

// Delete removes a row by primary key
func (t *Table) Delete(pk any) error {
	if t.PrimaryIndex == nil {
		return fmt.Errorf("table has no primary key index")
	}

	index, exists := t.PrimaryIndex[pk]
	if !exists {
		return fmt.Errorf("row with primary key %v not found", pk)
	}

	// Remove row from slice
	t.Rows = append(t.Rows[:index], t.Rows[index+1:]...)

	// Remove from primary index
	delete(t.PrimaryIndex, pk)

	// Update indexes for remaining rows
	for i := index; i < len(t.Rows); i++ {
		for _, col := range t.Columns {
			if col.IsPrimaryKey {
				t.PrimaryIndex[t.Rows[i].Data[col.Name]] = i
				break
			}
		}
	}

	return nil
}
