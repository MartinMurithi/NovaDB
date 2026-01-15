package parser

import (
	"reflect"
	"testing"
)

func TestParseSelectSimple(t *testing.T) {
	sql := "SELECT id, name FROM users"
	q, err := ParseSelect(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if q.Table != "users" {
		t.Fatalf("expected table 'users', got %s", q.Table)
	}

	expectedCols := []string{"id", "name"}
	if !reflect.DeepEqual(q.Columns, expectedCols) {
		t.Fatalf("expected columns %v, got %v", expectedCols, q.Columns)
	}

	if len(q.Filters) != 0 {
		t.Fatalf("expected no filters, got %v", q.Filters)
	}
}

func TestParseSelectWithWhereInt(t *testing.T) {
	sql := "SELECT id FROM users WHERE id = 1"
	q, err := ParseSelect(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if q.Table != "users" {
		t.Fatalf("expected table 'users', got %s", q.Table)
	}

	if len(q.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %v", q.Filters)
	}

	filter := q.Filters[0]
	if filter.Column != "id" || filter.Operator != "=" || filter.Value != 1 {
		t.Fatalf("unexpected filter: %v", filter)
	}
}

func TestParseSelectWithWhereString(t *testing.T) {
	sql := "SELECT name FROM users WHERE name = 'Alice'"
	q, err := ParseSelect(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	filter := q.Filters[0]
	if filter.Column != "name" || filter.Operator != "=" || filter.Value != "Alice" {
		t.Fatalf("unexpected filter: %v", filter)
	}
}

func TestParseSelectStar(t *testing.T) {
	sql := "SELECT * FROM orders"
	q, err := ParseSelect(sql)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(q.Columns) != 1 || q.Columns[0] != "*" {
		t.Fatalf("expected columns ['*'], got %v", q.Columns)
	}
}