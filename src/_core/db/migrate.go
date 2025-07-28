package db

import (
	"database/sql"
	"effective_mobile/src/_core/config"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(direction string) error {
	direction = strings.ToLower(direction)
	if direction != "up" && direction != "down" {
		return fmt.Errorf("invalid direction: %s (allowed: up/down)", direction)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error load config: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s search_path=%s sslmode=disable TimeZone=UTC",
		cfg.DB.Host,
		cfg.DB.User,
		url.QueryEscape(cfg.DB.Password),
		cfg.DB.Name,
		cfg.DB.Port,
		cfg.DB.Schema,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error connect DB: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(nil)
	goose.SetTableName("_migrations")

	switch direction {
	case "up":
		if err := goose.Up(db, "migrations"); err != nil {
			return fmt.Errorf("error applying migrations: %w", err)
		}
		log.Println("Migrations applied successfully!")
	case "down":
		if err := goose.Down(db, "migrations"); err != nil {
			return fmt.Errorf("error reverting migrations: %w", err)
		}
		log.Println("Migrations reverted successfully!")
	}

	return nil
}
