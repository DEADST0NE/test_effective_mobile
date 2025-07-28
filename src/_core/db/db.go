package db

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"effective_mobile/src/_core/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Error load config: %v", err)
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("Error connect DB")
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic("Error get underlying sql.DB")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Connected DB")

	return db, nil
}
