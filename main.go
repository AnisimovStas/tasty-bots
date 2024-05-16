package main

import (
	"context"
	"log"
	"tasty-bots/cmd"
	"tasty-bots/internal/storage/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlitePath = "data/sqlite/db.db"
)

func main() {

	s, err := sqlite.New(sqlitePath)
	if err != nil {
		log.Fatalf("can't connect to db %w", err)
	}

	err = s.Init(context.Background())
	if err != nil {
		log.Fatalf("can't init db %w", err)
	}

	cmd.Execute()
}
