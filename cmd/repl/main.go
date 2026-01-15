package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"

	"github.com/MartinMurithi/NovaDB.git/internal/engine"
	"github.com/MartinMurithi/NovaDB.git/internal/parser"
	"github.com/MartinMurithi/NovaDB.git/internal/planner"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

// SQL keywords for highlighting
var sqlKeywords = []string{
	"SELECT", "FROM", "WHERE", "INSERT", "INTO", "VALUES",
	"UPDATE", "SET", "DELETE", "AND", "OR",
}

// highlightSQL colors keywords
func highlightSQL(sql string) string {
	for _, kw := range sqlKeywords {
		sql = strings.ReplaceAll(sql, kw, "\033[1;34m"+kw+"\033[0m")
	}
	return sql
}

// PrintRows prints rows in table format, aligned and respecting column order
func PrintRows(rows []*storage.Row, columns []string, table *storage.Table) {
	if len(rows) == 0 {
		fmt.Println("(no rows)")
		return
	}

	// Determine columns
	if len(columns) == 1 && columns[0] == "*" {
		columns = []string{}
		for _, c := range table.Columns {
			columns = append(columns, c.Name)
		}
	}

	// Determine max width per column
	widths := make(map[string]int)
	for _, col := range columns {
		widths[col] = len(col)
	}
	for _, row := range rows {
		for _, col := range columns {
			valLen := len(fmt.Sprintf("%v", row.Data[col]))
			if valLen > widths[col] {
				widths[col] = valLen
			}
		}
	}

	// Print header
	for _, col := range columns {
		fmt.Printf("%-*s ", widths[col], col)
	}
	fmt.Println()
	for _, col := range columns {
		fmt.Printf("%s ", strings.Repeat("-", widths[col]))
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for _, col := range columns {
			val := row.Data[col]
			fmt.Printf("%-*v ", widths[col], val)
		}
		fmt.Println()
	}
}

func main() {
	// --------------------------
	// 1. Initialize DB and Engine
	// --------------------------
	db := storage.NewDatabase()
	eng := engine.NewEngine(db)

	// Create example table
	users, _ := db.CreateTable("users")
	users.AddColumn(&storage.Column{Name: "id", ColumnType: storage.IntType, IsPrimaryKey: true})
	users.AddColumn(&storage.Column{Name: "names", ColumnType: storage.TextType})
	users.AddColumn(&storage.Column{Name: "age", ColumnType: storage.IntType})

	// Insert initial data
	eng.Insert("users", map[string]any{"id": 1, "names": "Alice", "age": 30})
	eng.Insert("users", map[string]any{"id": 2, "names": "Bob", "age": 25})
	eng.Insert("users", map[string]any{"id": 3, "names": "Charlie", "age": 22})

	// --------------------------
	// 2. Setup REPL with readline
	// --------------------------
	rl, err := readline.New("> ")
	if err != nil {
		log.Fatalf("readline failed: %v", err)
	}
	defer rl.Close()

	fmt.Println("NovaDB REPL. Type 'exit;' to quit. End SQL with ';'")

	var buffer strings.Builder
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Collect multi-line SQL until ';'
		buffer.WriteString(" " + line)
		if !strings.HasSuffix(line, ";") {
			continue
		}

		sql := strings.TrimSpace(buffer.String())
		buffer.Reset()

		if strings.EqualFold(sql, "exit;") {
			fmt.Println("Bye!")
			break
		}

		// Highlight SQL
		fmt.Println(highlightSQL(sql))

		// --------------------------
		// Parse
		// --------------------------
		query, err := parser.Parse(sql)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			continue
		}

		// --------------------------
		// Plan
		// --------------------------
		plan, err := planner.CreatePlan(query)
		if err != nil {
			fmt.Printf("Planner error: %v\n", err)
			continue
		}

		// --------------------------
		// Execute
		// --------------------------
		rows, err := eng.ExecutePlan(plan)
		if err != nil {
			fmt.Printf("Execution error: %v\n", err)
			continue
		}

		// --------------------------
		// Print results
		// --------------------------
		table := eng.DB().Tables[plan.TableName]
		PrintRows(rows, plan.Columns, table)
	}
}