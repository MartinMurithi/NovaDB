package planner

import (
	"fmt"

	"github.com/MartinMurithi/NovaDB.git/internal/parser"
)

// --------------------------
// Plan Types
// --------------------------


type PlanType string

const (
	// DML
	SelectPlan PlanType = "SELECT"
	InsertPlan PlanType = "INSERT"
	UpdatePlan PlanType = "UPDATE"
	DeletePlan PlanType = "DELETE"

	// DDL
	CreateTablePlan   PlanType = "CREATE_TABLE"
	AddColumnPlan     PlanType = "ADD_COLUMN"
	ShowTablesPlan    PlanType = "SHOW_TABLES"
	DescribeTablePlan PlanType = "DESCRIBE_TABLE"
)


// Filter represents a WHERE clause condition
type Filter struct {
	Column   string
	Operator string // <, >, <=, >=, =, !=
	Value    any
}

// Plan represents an executable plan for a query
type Plan struct {
	Type      PlanType
	TableName string

	// SELECT
	Columns []string
	Filters []Filter

	// INSERT / UPDATE
	Values map[string]any

	// DDL
	ColumnsToAdd []string // For ADD COLUMN
	ColumnTypes  []string // Types for ADD COLUMN
}

// --------------------------
// CreatePlan
// --------------------------

func CreatePlan(q *parser.Query) (*Plan, error) {
	switch q.Type {

	// --------------------------
	case parser.CreateTableQuery:
		return &Plan{
			Type:      CreateTablePlan,
			TableName: q.Table,
		}, nil

	// --------------------------
	case parser.AddColumnQuery:
		return &Plan{
			Type:         AddColumnPlan,
			TableName:    q.Table,
			ColumnsToAdd: q.Columns,
			ColumnTypes:  q.ColumnTypes,
		}, nil

	// --------------------------
	case parser.ShowTablesQuery:
		return &Plan{
			Type: ShowTablesPlan,
		}, nil

	// --------------------------
	case parser.DescribeTableQuery:
		return &Plan{
			Type:      DescribeTablePlan,
			TableName: q.Table,
		}, nil

	// --------------------------
	case parser.SelectQuery:
		cols := q.Columns
		if len(cols) == 0 {
			cols = []string{"*"}
		}

		filters := make([]Filter, len(q.Filters))
		for i, f := range q.Filters {
			filters[i] = Filter{
				Column:   f.Column,
				Operator: f.Operator,
				Value:    f.Value,
			}
		}

		return &Plan{
			Type:      SelectPlan,
			TableName: q.Table,
			Columns:   cols,
			Filters:   filters,
		}, nil

	// --------------------------
	case parser.InsertQuery:
		values := make(map[string]any)
		for _, a := range q.Assignments {
			values[a.Column] = a.Value
		}

		return &Plan{
			Type:      InsertPlan,
			TableName: q.Table,
			Values:    values,
		}, nil

	// --------------------------
	case parser.UpdateQuery:
		values := make(map[string]any)
		for _, a := range q.Assignments {
			values[a.Column] = a.Value
		}

		filters := make([]Filter, len(q.Filters))
		for i, f := range q.Filters {
			filters[i] = Filter{
				Column:   f.Column,
				Operator: f.Operator,
				Value:    f.Value,
			}
		}

		return &Plan{
			Type:      UpdatePlan,
			TableName: q.Table,
			Values:    values,
			Filters:   filters,
		}, nil

	// --------------------------
	case parser.DeleteQuery:
		filters := make([]Filter, len(q.Filters))
		for i, f := range q.Filters {
			filters[i] = Filter{
				Column:   f.Column,
				Operator: f.Operator,
				Value:    f.Value,
			}
		}

		return &Plan{
			Type:      DeletePlan,
			TableName: q.Table,
			Filters:   filters,
		}, nil

	// --------------------------
	default:
		return nil, fmt.Errorf("unsupported query type %s", q.Type)
	}
}

