// package web

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/MartinMurithi/NovaDB.git/internal/engine"
// 	"github.com/MartinMurithi/NovaDB.git/internal/parser"
// 	"github.com/MartinMurithi/NovaDB.git/internal/planner"
// 	"github.com/MartinMurithi/NovaDB.git/internal/storage"
// 	"github.com/gin-gonic/gin"
// )

// func Run(db *storage.Database, eng *engine.Engine, addr string) {
// 	r := gin.Default()

// 	// --------------------------
// 	// Serve static files (UI)
// 	// --------------------------
// 	r.Static("/static", "/home/martin-wachira/Martin/NovaDB/internal/web/static") // folder with index.html, JS, CSS

// 	// Main page
// 	r.GET("/", func(c *gin.Context) {
// 		c.File("/home/martin-wachira/Martin/NovaDB/internal/web/static/index.html")
// 	})

// 	// SPA fallback (for client-side routing)
// 	r.NoRoute(func(c *gin.Context) {
// 		c.File("/home/martin-wachira/Martin/NovaDB/internal/web/static/index.html")
// 	})

// 	// --------------------------
// 	// Table management endpoints
// 	// --------------------------
// 	r.GET("/tables", func(c *gin.Context) {
// 		tables := []string{}
// 		for name := range db.Tables {
// 			tables = append(tables, name)
// 		}
// 		c.JSON(http.StatusOK, gin.H{"tables": tables})
// 	})

// 	r.POST("/table", func(c *gin.Context) {
// 		var body struct {
// 			Name string `json:"name"`
// 		}
// 		if err := c.BindJSON(&body); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		_, err := db.CreateTable(body.Name)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	r.GET("/table/:name/describe", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		table, ok := db.Tables[tableName]
// 		if !ok {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
// 			return
// 		}
// 		cols := []gin.H{}
// 		for _, col := range table.Columns {
// 			cols = append(cols, gin.H{"name": col.Name, "type": col.ColumnType})
// 		}
// 		c.JSON(http.StatusOK, gin.H{"columns": cols})
// 	})

// 	r.POST("/table/:name/column", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		table, ok := db.Tables[tableName]
// 		if !ok {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
// 			return
// 		}

// 		var body struct {
// 			Name string `json:"name"`
// 			Type string `json:"type"`
// 		}
// 		if err := c.BindJSON(&body); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		table.AddColumn(&storage.Column{
// 			Name:       body.Name,
// 			ColumnType: storage.ColumnType(body.Type),
// 		})

// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	// --------------------------
// 	// CRUD endpoints
// 	// --------------------------
// 	r.GET("/table/:name", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		sql := fmt.Sprintf("SELECT * FROM %s;", tableName)
// 		query, err := parser.Parse(sql)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		plan, _ := planner.CreatePlan(query)
// 		rows, err := eng.ExecutePlan(plan)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, rows)
// 	})

// 	r.POST("/table/:name", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		var body map[string]any
// 		if err := c.BindJSON(&body); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		cols := []string{}
// 		vals := []string{}
// 		for k, v := range body {
// 			cols = append(cols, k)
// 			switch val := v.(type) {
// 			case string:
// 				vals = append(vals, fmt.Sprintf("'%s'", val))
// 			default:
// 				vals = append(vals, fmt.Sprintf("%v", val))
// 			}
// 		}

// 		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
// 			tableName,
// 			strings.Join(cols, ", "),
// 			strings.Join(vals, ", "),
// 		)

// 		query, _ := parser.Parse(sql)
// 		plan, _ := planner.CreatePlan(query)
// 		_, err := eng.ExecutePlan(plan)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	r.PUT("/table/:name/:id", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		id := c.Param("id")
// 		var body map[string]any
// 		if err := c.BindJSON(&body); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		setParts := []string{}
// 		for k, v := range body {
// 			switch val := v.(type) {
// 			case string:
// 				setParts = append(setParts, fmt.Sprintf("%s='%s'", k, val))
// 			default:
// 				setParts = append(setParts, fmt.Sprintf("%s=%v", k, val))
// 			}
// 		}

// 		sql := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s;", tableName, strings.Join(setParts, ", "), id)
// 		query, _ := parser.Parse(sql)
// 		plan, _ := planner.CreatePlan(query)
// 		_, err := eng.ExecutePlan(plan)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	r.DELETE("/table/:name/:id", func(c *gin.Context) {
// 		tableName := c.Param("name")
// 		id := c.Param("id")
// 		sql := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", tableName, id)
// 		query, _ := parser.Parse(sql)
// 		plan, _ := planner.CreatePlan(query)
// 		_, err := eng.ExecutePlan(plan)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	// --------------------------
// 	// Start server
// 	// --------------------------
// 	log.Printf("NovaDB web UI running at http://%s", addr)
// 	if err := r.Run(addr); err != nil {
// 		log.Fatalf("failed to run server: %v", err)
// 	}
// }



package web

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

// Run starts the web server and API for NovaDB
func Run(db *storage.Database, eng *engine.Engine, addr string) {
	r := gin.Default()

	// --------------------------
	// Serve static files (UI)
	// --------------------------
	r.Static("/static", "/home/martin-wachira/Martin/NovaDB/internal/web/static")

	// Main page
	r.GET("/", func(c *gin.Context) {
		c.File("/home/martin-wachira/Martin/NovaDB/internal/web/static/index.html")
	})

	// SPA fallback for client-side routing
	r.NoRoute(func(c *gin.Context) {
		c.File("/home/martin-wachira/Martin/NovaDB/internal/web/static/index.html")
	})

	// --------------------------
	// Table management endpoints
	// --------------------------
	r.GET("/tables", func(c *gin.Context) {
		tables := []string{}
		for name := range db.Tables {
			tables = append(tables, name)
		}
		c.JSON(http.StatusOK, gin.H{"tables": tables})
	})

	r.POST("/table", func(c *gin.Context) {
		var body struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.CreateTable(body.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/table/:name/describe", func(c *gin.Context) {
		tableName := c.Param("name")
		table, ok := db.Tables[tableName]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}

		cols := []gin.H{}
		for _, col := range table.Columns {
			cols = append(cols, gin.H{"name": col.Name, "type": col.ColumnType})
		}
		c.JSON(http.StatusOK, gin.H{"columns": cols})
	})

	r.POST("/table/:name/column", func(c *gin.Context) {
		tableName := c.Param("name")
		table, ok := db.Tables[tableName]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}

		var body struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert UI type string to storage.ColumnType
		var colType storage.ColumnType
		switch strings.ToLower(body.Type) {
		case "int":
			colType = storage.IntType
		case "text":
			colType = storage.TextType
		case "float":
			colType = storage.FloatType
		case "bool":
			colType = storage.BoolType
		case "date":
			colType = storage.DateType
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unknown column type: %s", body.Type)})
			return
		}

		table.AddColumn(&storage.Column{
			Name:       body.Name,
			ColumnType: colType,
		})

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --------------------------
	// CRUD endpoints
	// --------------------------
	r.GET("/table/:name", func(c *gin.Context) {
		tableName := c.Param("name")
		sql := fmt.Sprintf("SELECT * FROM %s;", tableName)
		query, err := parser.Parse(sql)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		plan, _ := planner.CreatePlan(query)
		rows, err := eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, rows)
	})

	r.POST("/table/:name", func(c *gin.Context) {
		tableName := c.Param("name")
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

		sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
			tableName,
			strings.Join(cols, ", "),
			strings.Join(vals, ", "),
		)

		query, _ := parser.Parse(sql)
		plan, _ := planner.CreatePlan(query)
		_, err := eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.PUT("/table/:name/:id", func(c *gin.Context) {
		tableName := c.Param("name")
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

		sql := fmt.Sprintf("UPDATE %s SET %s WHERE id=%s;", tableName, strings.Join(setParts, ", "), id)
		query, _ := parser.Parse(sql)
		plan, _ := planner.CreatePlan(query)
		_, err := eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.DELETE("/table/:name/:id", func(c *gin.Context) {
		tableName := c.Param("name")
		id := c.Param("id")
		sql := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", tableName, id)
		query, _ := parser.Parse(sql)
		plan, _ := planner.CreatePlan(query)
		_, err := eng.ExecutePlan(plan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// --------------------------
	// Start server
	// --------------------------
	log.Printf("NovaDB web UI running at http://%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
