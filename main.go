package main

import (
	"log"
	"os"

	"finalyandex/pkg/db"
	"finalyandex/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	password := os.Getenv("TODO_PASSWORD")

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(password); err != nil {
		log.Fatal(err)
	}
}
