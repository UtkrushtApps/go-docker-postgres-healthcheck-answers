package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "github.com/lib/pq"
)

const maxRetries = 10
const retryInterval = 3 * time.Second

func main() {
	dbURL := getDBConnString()
	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err = db.PingContext(ctx)
			if err == nil {
				log.Println("Connected to database!")
				break
			}
		}
		log.Printf("Failed to connect to db (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
	}
	defer db.Close()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("DB unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("API server listening on :%s", getEnv("PORT", "8080"))
	http.ListenAndServe(":"+getEnv("PORT", "8080"), nil)
}

func getDBConnString() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "")
	dbname := getEnv("POSTGRES_DB", "postgres")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
