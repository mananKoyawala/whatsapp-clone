package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	// log.Printf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", username, password, host, port, dbname, sslmode)
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", username, password, host, port, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New("error during opening the database")
	}

	if err = db.Ping(); err != nil {
		return nil, err
		// + errors.New("ping error")
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
