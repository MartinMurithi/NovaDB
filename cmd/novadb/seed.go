package main

import (
	"log"

	"github.com/MartinMurithi/NovaDB.git/internal/engine"
	"github.com/MartinMurithi/NovaDB.git/internal/storage"
)

func Seed(db *storage.Database, eng *engine.Engine) {
	users, err := db.CreateTable("users")
	if err != nil {
		log.Println("seed: users table already exists")
		return
	}

	users.AddColumn(&storage.Column{
		Name:          "id",
		ColumnType:    storage.IntType,
		IsPrimaryKey:  true,
	})
	users.AddColumn(&storage.Column{
		Name:       "names",
		ColumnType: storage.TextType,
	})
	users.AddColumn(&storage.Column{
		Name:       "age",
		ColumnType: storage.IntType,
	})

	eng.Insert("users", map[string]any{
		"id":    1,
		"names": "Alice",
		"age":   30,
	})
	eng.Insert("users", map[string]any{
		"id":    2,
		"names": "Bob",
		"age":   25,
	})
}
