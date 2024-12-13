package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func ConnectToDB() (*sql.DB, error) {
	connString := os.Getenv("POSTGRES_CONNECTION_STRING")
	count := 1

	for {
		db, err := sql.Open("postgres", connString)
		if err != nil {
			log.Println("could not connect to postgres. retrying... ")
			count++
		} else {
			err = db.Ping()
			if err != nil {
				log.Println("postgres connection test failed. retrying...")
				count++
				db.Close()
			} else {
				return db, nil
			}
		}

		if count > 10 {
			return nil, err
		}

		log.Println("retrying in 1 second...")
		time.Sleep(1 * time.Second)
	}
}

func CreateTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
        )
    `
	_, err := db.Exec(query)
	return err
}

func UserExists(db *sql.DB, email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email= $1"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(db *sql.DB, email, password, first_name, last_name string, created_at time.Time) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (email, password, first_name, last_name, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.Exec(query, email, hashedPassword, first_name, last_name, created_at)
	return err
}

func hashPassword(plainPassword string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashBytes)
	return hash, nil
}
