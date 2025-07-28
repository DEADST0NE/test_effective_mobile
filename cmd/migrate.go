package main

import (
	"log"
	"os"

	"effective_mobile/src/_core/db"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/migrate/main.go <up|down>")
	}

	direction := os.Args[1]
	if err := db.RunMigrations(direction); err != nil {
		log.Fatal("Migration failed: ", err)
	}
}
