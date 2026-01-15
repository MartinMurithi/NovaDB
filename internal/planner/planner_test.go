package planner

import (
	"testing"

	"github.com/MartinMurithi/NovaDB.git/internal/parser"
)

func TestCreatePlan(t *testing.T) {
	q := &parser.Query{
		Type:    parser.SelectQuery,
		Table:   "users",
		Columns: []string{"id", "name"},
		Filters: []parser.Filter{
			{Column: "id", Operator: "=", Value: 1},
		},
	}

	plan, err := CreatePlan(q)
	if err != nil {
		t.Fatalf("planner failed: %v", err)
	}

	if plan.TableName != "users" {
		t.Fatalf("expected table 'users', got %s", plan.TableName)
	}

	if len(plan.Columns) != 2 || plan.Columns[0] != "id" || plan.Columns[1] != "name" {
		t.Fatalf("unexpected columns: %v", plan.Columns)
	}

	if len(plan.Filters) != 1 || plan.Filters[0].Column != "id" {
		t.Fatalf("unexpected filters: %v", plan.Filters)
	}
}