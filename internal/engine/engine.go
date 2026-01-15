package engine

import(
	// "github.com/MartinMurithi/NovaDB/internal/storage"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

type Engine struct {
	db *storage.Database
}

// NewEngine creates a new execution engine
func NewEngine(db *storage.Database) *Engine {
	return &Engine{
		db: db,
	}
}