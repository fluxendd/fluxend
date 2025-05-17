package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s port=5432",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_SSL_MODE"),
	)

	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	DB.DB.SetMaxOpenConns(10)
	DB.DB.SetMaxIdleConns(5)
	DB.DB.SetConnMaxLifetime(1 * time.Minute)
}

func GetDB() *sqlx.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return DB
}
