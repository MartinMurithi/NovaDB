package planner

import (
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/parser"
)

// Plan represents a simple execution plan
type Plan struct {
	TableName string
	Columns   []string
	Filters   []parser.Filter
}

// CreatePlan converts a parsed Query into an execution Plan
func CreatePlan(q *parser.Query) (*Plan, error) {
	if q.Type != parser.SelectQuery {
		return nil, fmt.Errorf("unsupported query type %s", q.Type)
	}

	// Columns check: if empty, default to "*"
	cols := q.Columns
	if len(cols) == 0 {
		cols = []string{"*"}
	}

	plan := &Plan{
		TableName: q.Table,
		Columns:   cols,
		Filters:   q.Filters,
	}

	return plan, nil
}