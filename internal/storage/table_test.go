package storage

import (
	"testing"
	"time"
)

func TestAddColumn(t *testing.T) {
	db := NewDatabase()

	table, err := db.CreateTable("users")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	if table.Name != "users" {
		t.Fatalf("expected table name 'users', got %v", table.Name)
	}

	if len(db.Tables) == 0 {
		t.Fatal("expected at least one table in database")
	}

	col := &Column{
		Name:         "id",
		ColumnType:   IntType,
		IsPrimaryKey: true,
		IsUnique:     true,
	}

	// Add column and check for errors
	if err := table.AddColumn(col); err != nil {
		t.Fatalf("failed to add column: %v", err)
	}

	// Verify column was added
	if len(table.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(table.Columns))
	}

	if table.Columns[0].Name != "id" {
		t.Fatalf("expected column name 'id', got %v", table.Columns[0].Name)
	}

	t.Logf("columns in table: %+v", table.Columns[0].Name)
}

func TestDropColumn(t *testing.T) {
	db := NewDatabase()

	// Create a table first
	table, err := db.CreateTable("users")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Add a couple of columns
	cols := []*Column{
		{Name: "id", ColumnType: IntType, IsPrimaryKey: true, IsUnique: true},
		{Name: "name", ColumnType: TextType},
	}

	for _, col := range cols {
		if err := table.AddColumn(col); err != nil {
			t.Fatalf("failed to add column %s: %v", col.Name, err)
		}
	}

	// Drop an existing column
	if err := table.DropColumn("id"); err != nil {
		t.Fatalf("failed to drop column 'id': %v", err)
	}

	// Verify column was removed
	for _, col := range table.Columns {
		if col.Name == "id" {
			t.Fatal("column 'id' still exists after DropColumn")
		}
	}

	// Ensure remaining columns are intact
	if len(table.Columns) != 1 || table.Columns[0].Name != "name" {
		t.Fatalf("expected remaining column 'name', got %+v", table.Columns)
	}

	t.Logf("columns after drop: %+v", table.Columns)

	// Try dropping a non-existent column
	err = table.DropColumn("nonexistent")
	if err == nil {
		t.Fatal("expected error when dropping non-existent column, got nil")
	}

	// Try dropping with empty name
	err = table.DropColumn("")
	if err == nil {
		t.Fatal("expected error when dropping column with empty name, got nil")
	}
}

func TestReadAllRows(t *testing.T) {
	db := NewDatabase()

	table, _ := db.CreateTable("products")

	cols := []*Column{
		{Name: "id", ColumnType: IntType, IsPrimaryKey: true, IsUnique: true},
		{Name: "name", ColumnType: TextType, IsPrimaryKey: false, IsUnique: true},
		{Name: "price", ColumnType: FloatType, IsPrimaryKey: false, IsUnique: false},
		{Name: "description", ColumnType: TextType, IsPrimaryKey: false, IsUnique: false},
		{Name: "is_discountable", ColumnType: BoolType, IsPrimaryKey: false, IsUnique: false},
		{Name: "added_date", ColumnType: DateType, IsPrimaryKey: false, IsUnique: false},
	}

	for _, col := range cols {
		if err := table.AddColumn(col); err != nil {
			t.Fatalf("an error occurred when adding columns %s", err)
		}
	}

	rows := []*Row{
		{Data: map[string]any{"id": 1, "name": "gaming mouse", "price": 1999.99, "description": "a white wireless gaming mouse", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 2, "name": "wireless keyboard", "price": 1500, "description": "a wireless gaming keyboard with neon LED lights", "is_discountable": true, "added_date": time.Now()}},
		{Data: map[string]any{"id": 3, "name": "oraimo Headphones", "price": 3000, "description": "bluetooth headphones", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 4, "name": "couch", "price": 10599.99, "description": "a grey L-Seat", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 5, "name": "electric table", "price": 20099.99, "description": "an electric height adjustable table", "is_discountable": true, "added_date": time.Now()}},
	}

	for _, row := range rows {
		if err := table.Insert(row); err != nil {
			t.Fatalf("an error occurred when adding row %v", err)
		}
	}

	all := table.GetRows()

	if len(all) != 5 {
		t.Fatalf("expected 5 rows but got %d", len(all))
	}

}
func TestReadRowByPK(t *testing.T) {
	db := NewDatabase()

	table, _ := db.CreateTable("products")

	cols := []*Column{
		{Name: "id", ColumnType: IntType, IsPrimaryKey: true, IsUnique: true},
		{Name: "name", ColumnType: TextType, IsPrimaryKey: false, IsUnique: true},
		{Name: "price", ColumnType: FloatType, IsPrimaryKey: false, IsUnique: false},
		{Name: "description", ColumnType: TextType, IsPrimaryKey: false, IsUnique: false},
		{Name: "is_discountable", ColumnType: BoolType, IsPrimaryKey: false, IsUnique: false},
		{Name: "added_date", ColumnType: DateType, IsPrimaryKey: false, IsUnique: false},
	}

	for _, col := range cols {
		if err := table.AddColumn(col); err != nil {
			t.Fatalf("an error occurred when adding columns %s", err)
		}
	}

	rows := []*Row{
		{Data: map[string]any{"id": 1, "name": "gaming mouse", "price": 1999.99, "description": "a white wireless gaming mouse", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 2, "name": "wireless keyboard", "price": 1500, "description": "a wireless gaming keyboard with neon LED lights", "is_discountable": true, "added_date": time.Now()}},
		{Data: map[string]any{"id": 3, "name": "oraimo Headphones", "price": 3000, "description": "bluetooth headphones", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 4, "name": "couch", "price": 10599.99, "description": "a grey L-Seat", "is_discountable": false, "added_date": time.Now()}},
		{Data: map[string]any{"id": 5, "name": "electric table", "price": 20099.99, "description": "an electric height adjustable table", "is_discountable": true, "added_date": time.Now()}},
	}

	for _, row := range rows {
		if err := table.Insert(row); err != nil {
			t.Fatalf("an error occurred when adding row %v", err)
		}
	}

	// Get row by primary key
	row, err := table.GetRowByPK(1)
	if err != nil || row.Data["name"] != "gaming mouse" {
		t.Fatal("failed to retrieve row by primary key")
	}

}

func TestReadRows(t *testing.T) {
	db := NewDatabase()
	table, _ := db.CreateTable("users")

	table.AddColumn(&Column{Name: "id", ColumnType: IntType, IsPrimaryKey: true})
	table.AddColumn(&Column{Name: "name", ColumnType: TextType})

	rows := []*Row{
		{Data: map[string]any{"id": 1, "name": "Alice"}},
		{Data: map[string]any{"id": 2, "name": "Bob"}},
	}

	for _, r := range rows {
		if err := table.Insert(r); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	// Filter rows
	filtered, err := table.FilterRows("name", "Bob")
	if err != nil || len(filtered) != 1 || filtered[0].Data["id"] != 2 {
		t.Fatal("failed to filter row by column")
	}
}

func TestUpdateRow(t *testing.T) {
	db := NewDatabase()
	table, _ := db.CreateTable("users")
	table.AddColumn(&Column{Name: "id", ColumnType: IntType, IsPrimaryKey: true})
	table.AddColumn(&Column{Name: "name", ColumnType: TextType})

	// Insert a row
	table.Insert(&Row{Data: map[string]any{"id": 1, "name": "Alice"}})

	// Update the row
	if err := table.Update(1, map[string]any{"name": "Alice Updated"}); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// Verify update
	row, err := table.GetRowByPK(1)
	if err != nil {
		t.Fatalf("failed to get row by primary key: %v", err)
	}

	if row.Data["name"] != "Alice Updated" {
		t.Fatal("row was not updated correctly")
	}

	t.Logf("updated row: %+v", row.Data)
}

func TestDeleteRow(t *testing.T) {
	db := NewDatabase()
	table, _ := db.CreateTable("users")
	table.AddColumn(&Column{Name: "id", ColumnType: IntType, IsPrimaryKey: true})
	table.AddColumn(&Column{Name: "name", ColumnType: TextType})

	// Insert a row
	table.Insert(&Row{Data: map[string]any{"id": 1, "name": "Alice"}})

	// Delete the row
	if err := table.Delete(1); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify deletion
	if len(table.Rows) != 0 {
		t.Fatal("row was not deleted")
	}

	_, err := table.GetRowByPK(1)
	if err == nil {
		t.Fatal("expected error retrieving deleted row, got nil")
	}

	t.Log("row deleted successfully")
}
