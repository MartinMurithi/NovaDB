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
