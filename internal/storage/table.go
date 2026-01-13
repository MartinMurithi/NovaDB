package storage

type Table struct {
	Name         string
	Columns      []*Column
	Rows         []*Row
	PrimaryIndex map[any]int  // Stores unique index for the primary key, Maps value of PK to row's index
}

// Insert, inserts a row into the table

// Validate types against columns
// Check PK uniqueness using primary Index
// Append row to row slices
// Update primary index

func (t *Table) Insert() error {

	return nil
}