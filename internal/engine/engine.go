package engine

import (
	// "github.com/MartinMurithi/NovaDB/internal/storage"
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/planner"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

type Engine struct {
	db *storage.Database
}

func NewEngine(db *storage.Database) *Engine {
	return &Engine{db: db}
}

func (e *Engine) DB() *storage.Database {
	return e.db
}

func (e *Engine) ExecutePlan(plan *planner.Plan) ([]*storage.Row, error) {
	switch plan.Type {

	// --------------------------
	case planner.CreateTablePlan:
		_, exists := e.db.Tables[plan.TableName]
		if exists {
			return nil, fmt.Errorf("table '%s' already exists", plan.TableName)
		}

		_, err := e.db.CreateTable(plan.TableName)
		if err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}

		return nil, nil // DDL commands return no rows

	// --------------------------
	case planner.AddColumnPlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}

		for i, col := range plan.ColumnsToAdd {
			colType := storage.TextType
			if i < len(plan.ColumnTypes) {
				colType = storage.ColumnType(plan.ColumnTypes[i])
			}
			t.AddColumn(&storage.Column{
				Name:       col,
				ColumnType: colType,
			})
		}

		return nil, nil

	// --------------------------
	case planner.ShowTablesPlan:
		rows := []*storage.Row{}
		for name := range e.db.Tables {
			rows = append(rows, &storage.Row{
				Data: map[string]any{"table_name": name},
			})
		}
		return rows, nil

	// --------------------------
	case planner.DescribeTablePlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}

		rows := []*storage.Row{}
		for _, col := range t.Columns {
			rows = append(rows, &storage.Row{
				Data: map[string]any{
					"name": col.Name,
					"type": col.ColumnType,
				},
			})
		}
		return rows, nil

	// --------------------------
	case planner.SelectPlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}
		return e.selectRows(plan, t)

	// --------------------------
	case planner.InsertPlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}
		return e.insertRow(plan, t)

	// --------------------------
	case planner.UpdatePlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}
		return e.updateRows(plan, t)

	// --------------------------
	case planner.DeletePlan:
		t, ok := e.db.Tables[plan.TableName]
		if !ok {
			return nil, fmt.Errorf("table '%s' does not exist", plan.TableName)
		}
		return e.deleteRows(plan, t)

	// --------------------------
	default:
		return nil, fmt.Errorf("unsupported plan type %s", plan.Type)
	}
}


// --------------------------
// SELECT helper
// --------------------------
func (e *Engine) selectRows(plan *planner.Plan, table *storage.Table) ([]*storage.Row, error) {
	rows := []*storage.Row{}

	for _, row := range table.Rows {
		if !matchesFilters(row, plan.Filters) {
			continue
		}

		// Project columns
		if len(plan.Columns) == 1 && plan.Columns[0] == "*" {
			rows = append(rows, row)
			continue
		}

		newData := make(map[string]any)
		for _, col := range plan.Columns {
			if val, ok := row.Data[col]; ok {
				newData[col] = val
			} else {
				return nil, fmt.Errorf("column '%s' does not exist in table '%s'", col, plan.TableName)
			}
		}
		rows = append(rows, &storage.Row{Data: newData})
	}

	return rows, nil
}

// --------------------------
// INSERT helper
// --------------------------
func (e *Engine) insertRow(plan *planner.Plan, table *storage.Table) ([]*storage.Row, error) {
	newRow := &storage.Row{Data: make(map[string]any)}
	for col, val := range plan.Values {
		if !e.TableHasColumn(plan.TableName, col) {
			return nil, fmt.Errorf("column '%s' does not exist in table '%s'", col, plan.TableName)
		}
		newRow.Data[col] = val
	}

	table.Rows = append(table.Rows, newRow)
	return []*storage.Row{newRow}, nil
}

// --------------------------
// UPDATE helper
// --------------------------
func (e *Engine) updateRows(plan *planner.Plan, table *storage.Table) ([]*storage.Row, error) {
	updated := []*storage.Row{}

	for _, row := range table.Rows {
		if matchesFilters(row, plan.Filters) {
			for col, val := range plan.Values {
				if !e.TableHasColumn(plan.TableName, col) {
					return nil, fmt.Errorf("column '%s' does not exist in table '%s'", col, plan.TableName)
				}
				row.Data[col] = val
			}
			updated = append(updated, row)
		}
	}

	return updated, nil
}

// --------------------------
// DELETE helper
// --------------------------
func (e *Engine) deleteRows(plan *planner.Plan, table *storage.Table) ([]*storage.Row, error) {
	remaining := []*storage.Row{}
	deleted := []*storage.Row{}

	for _, row := range table.Rows {
		if matchesFilters(row, plan.Filters) {
			deleted = append(deleted, row)
		} else {
			remaining = append(remaining, row)
		}
	}

	table.Rows = remaining
	return deleted, nil
}

// --------------------------
// Filters & column helpers
// --------------------------

func matchesFilters(row *storage.Row, filters []planner.Filter) bool {
	for _, f := range filters {
		val, ok := row.Data[f.Column]
		if !ok {
			return false
		}

		switch f.Operator {
		case "=":
			if val != f.Value {
				return false
			}
		case "!=":
			if val == f.Value {
				return false
			}
		case "<":
			if toFloat(val) >= toFloat(f.Value) {
				return false
			}
		case "<=":
			if toFloat(val) > toFloat(f.Value) {
				return false
			}
		case ">":
			if toFloat(val) <= toFloat(f.Value) {
				return false
			}
		case ">=":
			if toFloat(val) < toFloat(f.Value) {
				return false
			}
		default:
			panic("unsupported operator: " + f.Operator)
		}
	}
	return true
}

// Helper to convert int/float to float64 for comparison
func toFloat(v any) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case float64:
		return n
	default:
		panic(fmt.Sprintf("cannot convert %T to float64 for comparison", v))
	}
}

func (e *Engine) TableHasColumn(tableName, col string) bool {
	t, ok := e.db.Tables[tableName]
	if !ok {
		return false
	}
	for _, c := range t.Columns {
		if c.Name == col {
			return true
		}
	}
	return false
}
