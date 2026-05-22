package main

import (
	"flag"
	"log"
	"os"

	"crm-backend/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dir := flag.String("dir", "migrations", "migrations directory")
	flag.Parse()

	cfg := config.LoadConfig()
	db, err := goose.OpenDBWithDriver("pgx", cfg.DSN())
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"up"}
	}

	command := args[0]
	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("migrate status: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s (use: up | down | status)", command)
	}
	os.Exit(0)
}
