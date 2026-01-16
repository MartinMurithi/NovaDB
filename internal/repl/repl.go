package repl

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
// --------------------------
var sqlKeywords = []string{
	"SELECT", "FROM", "WHERE", "INSERT", "INTO", "VALUES",
	"UPDATE", "SET", "DELETE", "AND", "OR",
	"CREATE", "TABLE", "ALTER", "ADD", "COLUMN",
	"SHOW", "DESCRIBE",
}

func highlightSQL(sql string) string {
	for _, kw := range sqlKeywords {
		sql = strings.ReplaceAll(sql, kw, "\033[1;34m"+kw+"\033[0m")
	}
	return sql
}

// --------------------------
// PrintRows: nicely format output
// --------------------------
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

// --------------------------
// Main REPL
// --------------------------
func Run(db *storage.Database, eng *engine.Engine) {

	// --------------------------
	// 2. Setup REPL
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
		// Print results based on plan type
		// --------------------------
		switch plan.Type {
		case planner.SelectPlan:
			table, ok := db.Tables[plan.TableName]
			if !ok {
				fmt.Printf("Table %s not found\n", plan.TableName)
				continue
			}
			PrintRows(rows, plan.Columns, table)

		case planner.ShowTablesPlan:
			fmt.Println("Tables:")
			for _, r := range rows {
				fmt.Println(" -", r.Data["table_name"])
			}

		case planner.DescribeTablePlan:
			table, ok := db.Tables[plan.TableName]
			if !ok {
				fmt.Printf("Table %s not found\n", plan.TableName)
				continue
			}
			fmt.Printf("Columns in %s:\n", table.Name)
			for _, col := range table.Columns {
				fmt.Printf(" - %s (%s)\n", col.Name, col.ColumnType)
			}

		default:
			fmt.Printf("%s executed successfully\n", sql)
		}
	}
}