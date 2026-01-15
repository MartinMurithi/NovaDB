package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryType represents SQL query types
type QueryType string

const (
	SelectQuery QueryType = "SELECT"
)

// Filter represents a WHERE clause condition
type Filter struct {
	Column   string
	Operator string // Only "=" supported for now
	Value    any
}

// Query represents a parsed SQL query
type Query struct {
	Type    QueryType
	Table   string
	Columns []string
	Filters []Filter
}

// ParseSelect parses a simple SELECT SQL string
// Supported syntax:
// "SELECT col1, col2 FROM table"
// "SELECT * FROM table"
// "SELECT col1 FROM table WHERE col2 = 123"
// "SELECT col1 FROM table WHERE col2 = 'abc'"
func ParseSelect(sql string) (*Query, error) {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	if !strings.HasPrefix(sqlUpper, "SELECT") {
		return nil, fmt.Errorf("only SELECT queries supported")
	}

	parts := strings.SplitN(sql, "FROM", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}

	// Parse columns
	colsPart := strings.TrimSpace(parts[0][len("SELECT "):])
	columns := []string{}
	for _, c := range strings.Split(colsPart, ",") {
		columns = append(columns, strings.TrimSpace(c))
	}

	// Parse table and optional WHERE
	tablePart := strings.TrimSpace(parts[1])
	tableName := tablePart
	var filters []Filter

	if strings.Contains(strings.ToUpper(tablePart), "WHERE") {
		subParts := strings.SplitN(tablePart, "WHERE", 2)
		tableName = strings.TrimSpace(subParts[0])
		whereClause := strings.TrimSpace(subParts[1])

		// Only support simple equality: column = value
		parts := strings.SplitN(whereClause, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("only simple equality WHERE supported")
		}

		col := strings.TrimSpace(parts[0])
		valStr := strings.TrimSpace(parts[1])

		// Parse value as int if possible, else string
		var val any
		if i, err := strconv.Atoi(valStr); err == nil {
			val = i
		} else {
			val = strings.Trim(valStr, "'") // remove quotes for strings
		}

		filters = append(filters, Filter{
			Column:   col,
			Operator: "=",
			Value:    val,
		})
	}

	return &Query{
		Type:    SelectQuery,
		Table:   tableName,
		Columns: columns,
		Filters: filters,
	}, nil
}