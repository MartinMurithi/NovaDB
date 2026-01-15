package engine

import (
	"testing"

	"github.com/MartinMurithi/NovaDB.git/internal/parser"
	"github.com/MartinMurithi/NovaDB.git/internal/planner"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

func setupEngine() (*Engine, *storage.Table, error) {
	db := storage.NewDatabase()
	engine := NewEngine(db)

	table, _ := db.CreateTable("users")

	cols := []*storage.Column{
		{
			Name:         "id",
			ColumnType:   storage.IntType,
			IsPrimaryKey: true,
			IsUnique:     true,
		},
		{
			Name:         "names",
			ColumnType:   storage.TextType,
			IsPrimaryKey: false,
			IsUnique:     false,
		},
		{
			Name:         "email",
			ColumnType:   storage.TextType,
			IsPrimaryKey: false,
			IsUnique:     true,
		},
		{
			Name:         "username",
			ColumnType:   storage.TextType,
			IsPrimaryKey: false,
			IsUnique:     true,
		},
	}

	for _, col := range cols {
		if err := table.AddColumn(col); err != nil {
			return nil, nil, err
		}
	}

	return engine, table, nil
}

func TestEngineInsert(t *testing.T) {
	engine, table, err := setupEngine()

	if err := engine.Insert("users", map[string]any{"id": 1, "names": "Alice", "email": "alice@test.com", "username": "alice"}); err != nil {
		t.Fatalf("engine insert failed: %v", err)
	}

	row, err := engine.GetByPK("users", 1)
	if err != nil {
		t.Fatalf("GetByPK failed: %v", err)
	}

	if row.Data["names"] != "Alice" {
		t.Fatal("row data mismatch")
	}

	if len(table.Rows) != 1 {
		t.Fatal("expected 1 row after insert")
	}
}

func TestEngineSelectAll(t *testing.T) {
	engine, _, err := setupEngine()

	engine.Insert("users", map[string]any{"id": 1, "names": "Alice", "email": "alice@test.com", "username": "alice"})
	engine.Insert("users", map[string]any{"id": 2, "names": "Jane", "email": "jane@test.com", "username": "jane"})

	rows, err := engine.SelectAll("users")
	if err != nil {
		t.Fatalf("SelectAll failed: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}

func TestEngineSelectByColumnValue(t *testing.T) {
	engine, _, err := setupEngine()

	engine.Insert("users", map[string]any{"id": 1, "names": "Alice", "email": "alice@test.com", "username": "alice"})
	engine.Insert("users", map[string]any{"id": 2, "names": "Alice", "email": "ann@test.com", "username": "ann"})
	engine.Insert("users", map[string]any{"id": 3, "names": "Alex", "email": "alex@test.com", "username": "alex"})

	rows, err := engine.SelectByColumnValue("users", "email", "ann@test.com")
	if err != nil {
		t.Fatalf("SelectByColumnValue failed: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}

func TestEngineUpdate(t *testing.T) {
	engine, _, _ := setupEngine()

	engine.Insert("users", map[string]any{"id": 1, "names": "Alice"})

	if err := engine.Update("users", 1, map[string]any{"names": "Alice Updated"}); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	row, _ := engine.GetByPK("users", 1)
	if row.Data["names"] != "Alice Updated" {
		t.Fatal("row was not updated")
	}
}

func TestEngineDelete(t *testing.T) {
	engine, table, err := setupEngine()

	engine.Insert("users", map[string]any{"id": 1, "names": "Alice"})

	if err := engine.Delete("users", 1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if len(table.Rows) != 0 {
		t.Fatal("row was not deleted")
	}

	_, err = engine.GetByPK("users", 1)

	if err == nil {
		t.Fatal("expected GetByPK to fail after deletion")
	}
}

// ======== TESTS FOR QUERY AND PLAN EXECUTION ==========

func setupDB() (*storage.Database, *Engine) {
	db := storage.NewDatabase()
	eng := NewEngine(db)

	users, _ := db.CreateTable("users")
	users.AddColumn(&storage.Column{Name: "id", ColumnType: storage.IntType, IsPrimaryKey: true})
	users.AddColumn(&storage.Column{Name: "name", ColumnType: storage.TextType})

	eng.Insert("users", map[string]any{"id": 1, "name": "Alice"})
	eng.Insert("users", map[string]any{"id": 2, "name": "Bob"})
	eng.Insert("users", map[string]any{"id": 3, "name": "Charlie"})
	eng.Insert("users", map[string]any{"id": 4, "name": "Charles"})

	return db, eng
}
func TestExecutePlan_SelectAll(t *testing.T) {
	_, eng := setupDB()

	query, _ := parser.ParseSelect("SELECT * FROM users")
	plan, _ := planner.CreatePlan(query)

	rows, err := eng.ExecutePlan(plan)
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
}

func TestExecutePlan_SelectWithFilter(t *testing.T) {
	_, eng := setupDB()

	query, _ := parser.ParseSelect("SELECT id, name FROM users WHERE id = 2")
	plan, _ := planner.CreatePlan(query)

	rows, err := eng.ExecutePlan(plan)
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}

	if len(rows) != 1 || rows[0].Data["name"] != "Bob" {
		t.Fatal("filter did not return expected row")
	}
}

func TestExecutePlan_NonExistentTable(t *testing.T) {
	_, eng := setupDB()

	query, _ := parser.ParseSelect("SELECT * FROM unknown")
	plan, _ := planner.CreatePlan(query)

	_, err := eng.ExecutePlan(plan)
	if err == nil {
		t.Fatal("expected error for non-existent table")
	}
}

func TestExecutePlan_NonExistentColumn(t *testing.T) {
	_, eng := setupDB()

	query, _ := parser.ParseSelect("SELECT age FROM users")
	plan, _ := planner.CreatePlan(query)

	_, err := eng.ExecutePlan(plan)
	if err == nil {
		t.Fatal("expected error for non-existent column")
	}
}

func TestExecutePlan_NonExistentFilterColumn(t *testing.T) {
	_, eng := setupDB()

	query, _ := parser.ParseSelect("SELECT id FROM users WHERE age = 10")
	plan, _ := planner.CreatePlan(query)

	_, err := eng.ExecutePlan(plan)
	if err == nil {
		t.Fatal("expected error for non-existent filter column")
	}
}