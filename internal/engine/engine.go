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

// NewEngine creates a new execution engine
func NewEngine(db *storage.Database) *Engine {
	return &Engine{
		db: db,
	}
}

// ExecutePlan dispatches based on plan type
func (e *Engine) ExecutePlan(plan *planner.Plan) ([]*storage.Row, error) {
	switch plan.Type {
	case planner.SelectPlan:
		return e.executeSelect(plan)
	case planner.InsertPlan:
		return nil, e.executeInsert(plan)
	case planner.UpdatePlan:
		return nil, e.executeUpdate(plan)
	case planner.DeletePlan:
		return nil, e.executeDelete(plan)
	default:
		return nil, fmt.Errorf("unsupported plan type %s", plan.Type)
	}
}

func (e *Engine) ExecutePlan(plan *planner.Plan) ([]*storage.Row, error) {
	// 1. Lookup table
	rows, err := e.SelectAll(plan.TableName)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	// 2. Validate filter columns exist
	for _, f := range plan.Filters {
		if !e.TableHasColumn(plan.TableName, f.Column) {
			return nil, fmt.Errorf("execution error: filter column '%s' does not exist in table '%s'", f.Column, plan.TableName)
		}
	}

	// 3. Apply filters
	if len(plan.Filters) > 0 {
		filtered := []*storage.Row{}
		for _, row := range rows {
			matches := true
			for _, f := range plan.Filters {
				rowVal, ok := row.Data[f.Column]
				if !ok {
					return nil, fmt.Errorf("execution error: row missing column '%s'", f.Column)
				}
				if rowVal != f.Value {
					matches = false
					break
				}
			}
			if matches {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}

	// 4. Validate requested columns exist
	if !(len(plan.Columns) == 1 && plan.Columns[0] == "*") {
		for _, col := range plan.Columns {
			if !e.TableHasColumn(plan.TableName, col) {
				return nil, fmt.Errorf("execution error: requested column '%s' does not exist in table '%s'", col, plan.TableName)
			}
		}
	}

	// 5. Project columns
	if !(len(plan.Columns) == 1 && plan.Columns[0] == "*") {
		for _, row := range rows {
			newVals := map[string]any{}
			for _, col := range plan.Columns {
				newVals[col] = row.Data[col]
			}
			row.Data = newVals
		}
	}

	return rows, nil
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
