package planner

import (
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/parser"
)

type PlanType string

const (
	SelectPlan PlanType = "SELECT"
	InsertPlan PlanType = "INSERT"
	UpdatePlan PlanType = "UPDATE"
	DeletePlan PlanType = "DELETE"
)

type Plan struct {
	Type      PlanType
	TableName string

	// SELECT
	Columns []string
	Filters []Filter

	// INSERT / UPDATE
	Values map[string]any
}


// CreatePlan converts a parsed Query into an execution Plan
func CreatePlan(q *parser.Query) (*Plan, error) {
	switch q.Type {

	case parser.SelectQuery:
		cols := q.Columns
		if len(cols) == 0 {
			cols = []string{"*"}
		}

		return &Plan{
			Type:      SelectPlan,
			TableName: q.Table,
			Columns:   cols,
			Filters:   q.Filters,
		}, nil

	case parser.InsertQuery:
		return &Plan{
			Type:      InsertPlan,
			TableName: q.Table,
			Values:    q.Values,
		}, nil

	case parser.UpdateQuery:
		return &Plan{
			Type:      UpdatePlan,
			TableName: q.Table,
			Filters:   q.Filters,
			Values:    q.Values,
		}, nil

	case parser.DeleteQuery:
		return &Plan{
			Type:      DeletePlan,
			TableName: q.Table,
			Filters:   q.Filters,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported query type %s", q.Type)
	}
}
