package parser

import (
	"strconv"
	"strings"
)


func parseValue(v string) any {
	v = strings.TrimSpace(v)

	if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
		return strings.Trim(v, "'")
	}

	if i, err := strconv.Atoi(v); err == nil {
		return i
	}

	return v
}

func parseWhere(whereClause string) []Filter {
	conditions := strings.Split(whereClause, "AND") // simple AND only
	filters := []Filter{}

	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)

		// Detect operator
		var op string
		for _, o := range []string{"!=", ">=", "<=", "=", ">", "<"} {
			if strings.Contains(cond, o) {
				op = o
				break
			}
		}

		if op == "" {
			panic("unsupported operator in WHERE")
		}

		parts := strings.SplitN(cond, op, 2)
		col := strings.TrimSpace(parts[0])
		valRaw := strings.TrimSpace(parts[1])

		// Remove quotes for string literals
		var val any
		if strings.HasPrefix(valRaw, "'") && strings.HasSuffix(valRaw, "'") {
			val = strings.Trim(valRaw, "'")
		} else {
			val = parseValue(valRaw) // parse numbers
		}

		filters = append(filters, Filter{
			Column:   col,
			Operator: op,
			Value:    val,
		})
	}

	return filters
}