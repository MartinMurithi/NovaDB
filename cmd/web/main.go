package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MartinMurithi/NovaDB.git/internal/engine"
	"github.com/MartinMurithi/NovaDB.git/internal/parser"
	"github.com/MartinMurithi/NovaDB.git/internal/planner"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// --------------------------
	// 1. Initialize DB and Engine
	// --------------------------
	db := storage.NewDatabase()
	eng := engine.NewEngine(db)

	// Example table "users"
	users, _ := db.CreateTable("users")
	users.AddColumn(&storage.Column{Name: "id", ColumnType: storage.IntType, IsPrimaryKey: true})
	users.AddColumn(&storage.Column{Name: "names", ColumnType: storage.TextType})
	users.AddColumn(&storage.Column{Name: "age", ColumnType: storage.IntType})

	// Insert sample data
	eng.Insert("users", map[string]any{"id": 1, "names": "Alice", "age": 30})
	eng.Insert("users", map[string]any{"id": 2, "names": "Bob", "age": 25})
	eng.Insert("users", map[string]any{"id": 3, "names": "Charlie", "age": 22})

	// --------------------------
	// 2. Setup Gin router
	// --------------------------
	r := gin.Default()

	// Generic GET /table/:name → SELECT * FROM table
	r.GET("/table/:name", func(c *gin.Context) {
		table := c.Param("name")
		sql := fmt.Sprintf("SELECT * FROM %s;", table)

		query, err := parser.Parse(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		plan, err := planner.CreatePlan(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rows, err := eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rows)
	})

	// Generic POST /table/:name → INSERT
	r.POST("/table/:name", func(c *gin.Context) {
		table := c.Param("name")
		var body map[string]any
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cols := []string{}
		vals := []string{}
		for k, v := range body {
			cols = append(cols, k)
			switch val := v.(type) {
			case string:
				vals = append(vals, fmt.Sprintf("'%s'", val))
			default:
				vals = append(vals, fmt.Sprintf("%v", val))
			}
		}

		sql := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s);",
			table,
			strings.Join(cols, ", "),
			strings.Join(vals, ", "),
		)

		query, err := parser.Parse(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		plan, err := planner.CreatePlan(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Generic PUT /table/:name/:id → UPDATE by primary key
	r.PUT("/table/:name/:id", func(c *gin.Context) {
		table := c.Param("name")
		id := c.Param("id")

		var body map[string]any
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setParts := []string{}
		for k, v := range body {
			switch val := v.(type) {
			case string:
				setParts = append(setParts, fmt.Sprintf("%s='%s'", k, val))
			default:
				setParts = append(setParts, fmt.Sprintf("%s=%v", k, val))
			}
		}

		sql := fmt.Sprintf(
			"UPDATE %s SET %s WHERE id=%s;",
			table,
			strings.Join(setParts, ", "),
			id,
		)

		query, err := parser.Parse(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		plan, err := planner.CreatePlan(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Generic DELETE /table/:name/:id → DELETE by primary key
	r.DELETE("/table/:name/:id", func(c *gin.Context) {
		table := c.Param("name")
		id := c.Param("id")

		sql := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", table, id)

		query, err := parser.Parse(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		plan, err := planner.CreatePlan(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err = eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --------------------------
	// 3. Run server
	// --------------------------
	fmt.Println("NovaDB web server running at http://localhost:7070")
	if err := r.Run(":7070"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
