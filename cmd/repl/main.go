package main

import (
	"fmt"
	"log"

	"github.com/MartinMurithi/NovaDB.git/internal/engine"
	"github.com/MartinMurithi/NovaDB.git/internal/parser"
	"github.com/MartinMurithi/NovaDB.git/internal/planner"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

func main(){

	// 1. Initialize database and engine
	db := storage.NewDatabase()
	eng := engine.NewEngine(db)

	// 2. Create table and columns
	users, _ := db.CreateTable("users")
	users.AddColumn(&storage.Column{Name: "id", ColumnType: storage.IntType, IsPrimaryKey: true})
	users.AddColumn(&storage.Column{Name: "name", ColumnType: storage.TextType})

	// 3. Insert rows
	eng.Insert("users", map[string]any{"id": 1, "name": "Alice"})
	eng.Insert("users", map[string]any{"id": 2, "name": "Bob"})
	eng.Insert("users", map[string]any{"id": 3, "name": "Charlie"})

	// 4. Example SQL string
	// sql := "SELECT id, name FROM users WHERE id = 2"
	sql1 := "SELECT * FROM users"

	// 5. Parse SQL
	query, err := parser.ParseSelect(sql1)
	if err != nil {
		log.Fatalf("parse failed: %v", err)
	}

	// 6. Plan query
	plan, err := planner.CreatePlan(query)
	if err != nil {
		log.Fatalf("planner failed: %v", err)
	}

	// 7. Execute plan
	rows, err := eng.ExecutePlan(plan)
	if err != nil {
		log.Fatalf("execution failed: %v", err)
	}

	// 8. Print results
	for _, row := range rows {
		fmt.Println(row.Data)
	}

}