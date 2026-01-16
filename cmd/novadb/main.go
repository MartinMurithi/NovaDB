package main

import (
	"flag"
	"log"

	"github.com/MartinMurithi/NovaDB.git/internal/engine"
	"github.com/MartinMurithi/NovaDB.git/internal/repl"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
	"github.com/MartinMurithi/NovaDB.git/internal/web"
)

func main() {
	mode := flag.String("mode", "repl", "repl | web | both")
	addr := flag.String("addr", ":7070", "http address")
	flag.Parse()

	db := storage.NewDatabase()
	eng := engine.NewEngine(db)

	Seed(db, eng)

	switch *mode {
	case "repl":
		repl.Run(db, eng)

	case "web":
		web.Run(db, eng, *addr)

	case "both":
		go web.Run(db, eng, *addr)
		repl.Run(db, eng)

	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
}
