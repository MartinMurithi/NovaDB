package parser

import (
	"fmt"
	"strings"
)

type QueryType string

const (
	SelectQuery QueryType = "SELECT"
	InsertQuery QueryType = "INSERT"
	UpdateQuery QueryType = "UPDATE"
	DeleteQuery QueryType = "DELETE"
)

type Filter struct {
	Column   string
	Operator string
	Value    any
}

type Assignment struct {
	Column string
	Value  any
}

type Query struct {
	Type  QueryType
	Table string

	// SELECT
	Columns []string

	// WHERE (shared)
	Filters []Filter

	// INSERT / UPDATE
	Assignments []Assignment
}

func Parse(sql string) (*Query, error) {
	sql = strings.TrimSpace(sql)
	sql = strings.TrimSuffix(sql, ";")
	sql = strings.ToUpper(sql[:6]) + sql[6:]

	switch {
	case strings.HasPrefix(sql, "SELECT"):
		return parseSelect(sql)
	case strings.HasPrefix(sql, "INSERT"):
		return parseInsert(sql)
	case strings.HasPrefix(sql, "UPDATE"):
		return parseUpdate(sql)
	case strings.HasPrefix(sql, "DELETE"):
		return parseDelete(sql)
	default:
		return nil, fmt.Errorf("unsupported SQL statement")
	}
}

func parseSelect(sql string) (*Query, error) {
	// SELECT a,b FROM table WHERE c=1
	parts := strings.Split(sql, "FROM")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}

	colsPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "SELECT"))
	columns := strings.Split(colsPart, ",")

	for i := range columns {
		columns[i] = strings.TrimSpace(columns[i])
	}

	rest := strings.TrimSpace(parts[1])
	tableAndWhere := strings.Split(rest, "WHERE")

	q := &Query{
		Type:    SelectQuery,
		Table:   strings.TrimSpace(tableAndWhere[0]),
		Columns: columns,
	}

	if len(tableAndWhere) == 2 {
		q.Filters = parseWhere(tableAndWhere[1])
	}

	return q, nil
}

func parseInsert(sql string) (*Query, error) {
	// INSERT INTO t (a,b) VALUES (1,2)
	parts := strings.Split(sql, "VALUES")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid INSERT syntax")
	}

	head := parts[0]
	valuesPart := parts[1]

	table := strings.Fields(head)[2]

	cols := strings.Split(
		strings.TrimSpace(head[strings.Index(head, "(")+1:strings.Index(head, ")")]),
		",",
	)

	vals := strings.Split(
		strings.Trim(valuesPart, " ()"),
		",",
	)

	if len(cols) != len(vals) {
		return nil, fmt.Errorf("columns/values mismatch")
	}

	assignments := []Assignment{}
	for i := range cols {
		assignments = append(assignments, Assignment{
			Column: strings.TrimSpace(cols[i]),
			Value:  parseValue(vals[i]),
		})
	}

	return &Query{
		Type:        InsertQuery,
		Table:       table,
		Assignments: assignments,
	}, nil
}

func parseUpdate(sql string) (*Query, error) {
	// UPDATE t SET a=1 WHERE id=2
	parts := strings.Split(sql, "SET")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid UPDATE syntax")
	}

	table := strings.Fields(parts[0])[1]
	setAndWhere := strings.Split(parts[1], "WHERE")

	assignments := []Assignment{}
	for _, a := range strings.Split(setAndWhere[0], ",") {
		kv := strings.Split(a, "=")
		assignments = append(assignments, Assignment{
			Column: strings.TrimSpace(kv[0]),
			Value:  parseValue(kv[1]),
		})
	}

	q := &Query{
		Type:        UpdateQuery,
		Table:       table,
		Assignments: assignments,
	}

	if len(setAndWhere) == 2 {
		q.Filters = parseWhere(setAndWhere[1])
	}

	return q, nil
}

func parseDelete(sql string) (*Query, error) {
	// DELETE FROM t WHERE id=1
	parts := strings.Split(sql, "WHERE")
	table := strings.Fields(parts[0])[2]

	q := &Query{
		Type:  DeleteQuery,
		Table: table,
	}

	if len(parts) == 2 {
		q.Filters = parseWhere(parts[1])
	}

	return q, nil
}
